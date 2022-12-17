package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
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

//go run main.go -container-path=/vagrant/sample/janus/assets/02.mp4 --player=2
func main() {

	containerPath := ""
	jaunsURL := ""
	player := ""
	flag.StringVar(&containerPath, "container-path", "", "path to the media file you want to playback")
	flag.StringVar(&jaunsURL, "janus-url", "ws://localhost:8188/janus", "ws://localhost:8188/janus")
	flag.StringVar(&player, "player", "123", "123")
	flag.Parse()

	if containerPath == "" {
		panic("-container-path must be specified")
	}


	peerConnectionConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
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
	//audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "audio", "pion")
	audioTrack, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "synced-video", "synced-video")
    if err != nil {
		panic(err)
	} else if _, err = peerConnection.AddTrack(audioTrack); err != nil {
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
    pipelineStr := fmt.Sprintf("filesrc location=\"%s\" ! decodebin ! audioresample ! audioconvert ! opusenc ! appsink name=audio", containerPath)
	pipeline = gst.CreatePipeline(pipelineStr, audioTrack, videoTrack)
	pipeline.Start()

    // Create Janus
	fmt.Println(jaunsURL)
	gateway, err := janus.Connect(jaunsURL)
	if err != nil {
		panic(err)
	}

	session, err := gateway.Create()
	if err != nil {
		panic(err)
	}

	handle, err := session.Attch("janus.plugin.audiobridge")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			if _, keepAliveErr := session.KeepAlive(); keepAliveErr != nil {
				panic(keepAliveErr)
			}

			time.Sleep(50 * time.Second)
		}
	}()

	go watchHandle(handle)

	roomId := 2532488013468683936

    _, err = handle.Request(map[string]interface{}{
        "request": "create",
		"room": roomId,
        "description": "audio bridge Room",
		"secret": "adminpwd",
		"sampling_rate": 16000,
		"spatial_audio": true,
		"record": false,
    })
	if err != nil {
		panic(err)
	}

	_, err = handle.Request(map[string]interface{}{
		"request": "rtp_forward",
		"room": roomId,
		"host": "192.168.56.168",
        "secret": "adminpwd",
		"port": 5012,
	})
	if err != nil {
		panic(err)
	}

	x3 := rand.NewSource(5)
    y3 := rand.New(x3)
	num, _ := strconv.Atoi(player)
	joinId := y3.Intn(200) + num
	msg, err := handle.Message(map[string]interface{}{
		"request": 	"join",
		"room":    	roomId,
		"id": 		joinId,
		"display": 	containerPath,
		"muted": 	true,
		"quality": 	4,
		"volume": 	100,
		"spatial_position": 100,
	}, nil)
	if err != nil {
		panic(err)
	}
	log.Println("audiobridge join: ", msg)


	msg, err = handle.Message(map[string]interface{}{
		"request": "configure",
		"muted":   false,
	}, map[string]interface{}{
		"type":    "offer",
		"sdp":     peerConnection.LocalDescription().SDP,
	})
	if err != nil {
		panic(err)
	}
	log.Println("audiobridge configure: ", msg)

	if msg.Jsep != nil {
		err = peerConnection.SetRemoteDescription(webrtc.SessionDescription{
			Type: webrtc.SDPTypeAnswer,
			SDP:  msg.Jsep["sdp"].(string),
		})
		if err != nil {
			panic(err)
		}
        pipeline.Play()
	}

	select {}

}