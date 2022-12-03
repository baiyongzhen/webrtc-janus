package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"example.com/webrtc-janus/pkg/gst"
	"example.com/webrtc-janus/pkg/janus"
	"github.com/pion/webrtc/v3"
)


func watchHandle(handle *janus.Handle) {
	// wait for event
	for {
		msg := <-handle.Events
		switch msg := msg.(type) {
		case *janus.SlowLinkMsg:
			log.Println("SlowLinkMsg type ", handle.ID)
		case *janus.MediaMsg:
			log.Println("MediaEvent type", msg.Type, " receiving ", msg.Receiving)
		case *janus.WebRTCUpMsg:
			log.Println("WebRTCUp type ", handle.ID)
		case *janus.HangupMsg:
			log.Println("HangupEvent type ", handle.ID)
		case *janus.EventMsg:
			log.Printf("EventMsg %+v", msg.Plugindata.Data)
		}
	}
}

// go run main.go -container-path=/vagrant/projects/webrtc-janus/assets/01.mp4
func main() {

	containerPath := ""
	jaunsURL := ""
	flag.StringVar(&containerPath, "container-path", "", "path to the media file you want to playback")
	flag.StringVar(&jaunsURL, "janus-url", "ws://localhost:8188/janus", "ws://localhost:8188/janus")
	flag.Parse()

	if containerPath == "" {
		panic("-container-path must be specified")
	}


	peerConnectionConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
			/*{
				URLs:       []string{"turn:numb.viagenie.ca"},
				Username:   "webrtc@live.com",
				Credential: "muazkh",
			},*/
		},
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
    }
	audioTrack := &webrtc.TrackLocalStaticSample{}
	videoTrack := &webrtc.TrackLocalStaticSample{}
	pipeline   := &gst.Pipeline{}

    // local webrtc peer connection
    peerConnection, err := webrtc.NewPeerConnection(peerConnectionConfig)
    if err != nil {
		panic(err)
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	// Create a audio track
	audioTrack, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "synced-video", "synced-video")
    if err != nil {
		panic(err)
	} else if _, err = peerConnection.AddTrack(audioTrack); err != nil {
		panic(err)
	}

	// Create a video track
	videoTrack, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/h264"}, "synced-video", "synced-video")
    if err != nil {
		panic(err)
	} else if _, err = peerConnection.AddTrack(videoTrack); err != nil {
		panic(err)
	}

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	if err = peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

    // Create gstreamer pipeline
	pipelineStr := fmt.Sprintf("filesrc location=\"%s\" ! decodebin name=demux ! queue ! x264enc bframes=0 speed-preset=veryfast key-int-max=60 ! video/x-h264,stream-format=byte-stream ! appsink name=video demux. ! queue ! audioconvert ! audioresample ! opusenc ! appsink name=audio", containerPath)
	pipeline = gst.CreatePipeline(pipelineStr, audioTrack, videoTrack)
	pipeline.Start()

    // Create Janus
	//gateway, err := janus.Connect("ws://localhost:8188/janus")
	fmt.Println(jaunsURL)
	gateway, err := janus.Connect(jaunsURL)
	if err != nil {
		panic(err)
	}

	session, err := gateway.Create()
	if err != nil {
		panic(err)
	}

	handle, err := session.Attch("janus.plugin.videoroom")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			if _, keepAliveErr := session.KeepAlive(); keepAliveErr != nil {
				panic(keepAliveErr)
			}

			time.Sleep(5 * time.Second)
		}
	}()

	go watchHandle(handle)

	roomId := 2532488013468683936

    _, err = handle.Request(map[string]interface{}{
        "room": roomId,
        "request": "create",
        "description": "Demo Room",
		"admin_key": "jaster",
        "secret": "adminpwd",
        "publishers": 6,
        "bitrate": 128000,
        "fir_freq": 10,
        "record": false,
        "videocodec": "h264",
    })
	if err != nil {
		panic(err)
	}

	joinmsg, err := handle.Message(map[string]interface{}{
		"request": "join",
		"ptype":   "publisher",
		"room":    roomId,
		"display": "test",
	}, nil)
	if err != nil {
		panic(err)
	}

    feedId := joinmsg.Plugindata.Data["id"]
    fmt.Printf("FeedID:= %.0f \n", feedId)

	msg, err := handle.Message(map[string]interface{}{
		"request": "publish",
		"audio":   true,
		"video":   true,
		"data":    false,
		"bitrate": 128000,
		"bitrate_cap": true,
	}, map[string]interface{}{
		"type":    "offer",
		"sdp":     peerConnection.LocalDescription().SDP,
		"trickle": false,
	})
	if err != nil {
		panic(err)
	}

	if msg.Jsep != nil {
		err = peerConnection.SetRemoteDescription(webrtc.SessionDescription{
			Type: webrtc.SDPTypeAnswer,
			SDP:  msg.Jsep["sdp"].(string),
		})
		if err != nil {
			panic(err)
		}

		// Start pushing buffers on these tracks
        pipeline.Play()
	}

	/*
	//request rtp_forward  in v1.x.x version
	streams := []*map[string]interface{}{
		{
			"mid": "0",
			"port": 5006,
		},
		{
			"mid": "1",
			"port": 5011,
		},
	}

	time.Sleep(8 * time.Second)
	_, err = handle.Message(map[string]interface{}{
		"request": "rtp_forward",
		"room": roomId,
		"publisher_id": feedId,
		"host": "192.168.56.168",
		"secret": "adminpwd",
		"streams": streams,
	}, nil)
	if err != nil {
		//panic(err)
	}
	*/

	_, err = handle.Message(map[string]interface{}{
		"request": "rtp_forward",
		"room": roomId,
		"publisher_id": feedId,
		"host": "192.168.56.168",
		//"host_family": "ipv4",
		"admin_key": "janusoverlord",
        "secret": "adminpwd",
		"audio_port": 5006,
		"video_port": 5011,
	}, nil)
	if err != nil {
		fmt.Println(err.Error())
		//panic(err)
	}




	select {}

}