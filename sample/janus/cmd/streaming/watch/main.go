package main

import (
	"fmt"
	"log"
	"time"

	"example.com/webrtc-janus/pkg/janus"
	"github.com/pion/webrtc/v3"
	_ "github.com/pion/webrtc/v3/pkg/media/oggwriter"
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


func main() {

	mediaEngine := &webrtc.MediaEngine{}

	// Setup the codecs you want to use.
	// Only support VP8 and OPUS, this makes our WebM muxer code simpler
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: "video/h264", ClockRate: 90000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
		PayloadType:        96,
	}, webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}
	if err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: "audio/opus", ClockRate: 48000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
		PayloadType:        111,
	}, webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
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

	// Create a new RTCPeerConnection
	var peerConnection *webrtc.PeerConnection
	peerConnection, err := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine)).NewPeerConnection(peerConnectionConfig)
	if err != nil {
		panic(err)
	}

	/*
	// Allow us to receive 1 audio track, and 1 video track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	} else if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}
	*/

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		fmt.Printf("Track has started, of type %d: %s \n", track.PayloadType(), track.Codec().RTPCodecCapability.MimeType)
		/*
		codec := track.Codec()
		fmt.Println(codec.MimeType)
		if codec.Name == webrtc.Opus {
			fmt.Println("Got Opus track, saving to disk as output.ogg")
			i, oggNewErr := oggwriter.New("output.ogg", codec.ClockRate, codec.Channels)
			if oggNewErr != nil {
				panic(oggNewErr)
			}
			//saveToDisk(i, track)
		} else if codec.Name == webrtc.VP8 {
			fmt.Println("Got VP8 track, saving to disk as output.ivf")
			i, ivfNewErr := ivfwriter.New("output.ivf")
			if ivfNewErr != nil {
				panic(ivfNewErr)
			}
			//saveToDisk(i, track)
		}
		*/
	})


	jaunsURL := "ws://localhost:8188/janus"

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

	handle, err := session.Attch("janus.plugin.streaming")
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

	// Get streaming list
	_, err = handle.Request(map[string]interface{}{
		"request": "list",
	})
	if err != nil {
		panic(err)
	}

	// Watch the second stream
	msg, err := handle.Message(map[string]interface{}{
		"request": "watch",
		"id":      6,
	}, nil)
	if err != nil {
		panic(err)
	}

	// Wait for the offer to be pasted
	if msg.Jsep != nil {
		err = peerConnection.SetRemoteDescription(webrtc.SessionDescription{
			Type: webrtc.SDPTypeOffer,
			SDP:  msg.Jsep["sdp"].(string),
		})
		if err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	//fmt.Println(answer.SDP)
	_, err = handle.Message(map[string]interface{}{
		"request": "start",
		"id": 6,
	}, map[string]interface{}{
		"type": "answer",
		"sdp": answer.SDP,
	})
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete


	select {}


}
