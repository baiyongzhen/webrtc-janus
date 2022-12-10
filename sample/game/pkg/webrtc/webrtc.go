// credit to https://github.com/poi5305/go-yuv2webRTC/blob/master/webrtc/webrtc.go
package webrtc

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"time"
	"runtime/debug"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/rs/xid"
	"example.com/webrtc-game/pkg/janus"
)


type InputDataPair struct {
	data int
	time time.Time
}

type WebFrame struct {
	Data      []byte
	Timestamp uint32
}

// Game Meta
type GameMeta struct {
	PlayerIndex int
}

type OnIceCallback func(candidate string)

// Encode encodes the input in base64
func Encode(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// Decode decodes the input from base64
func Decode(in string, obj interface{}) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		return err
	}

	return nil
}

// WebRTC connection
type WebRTC struct {
	ID string

	//이것은 조금 생각을 해봐야함.
	//connection  *webrtc.PeerConnection
	isConnected bool
	isClosed    bool
	// for yuvI420 image
	ImageChannel    chan WebFrame
	AudioChannel    chan []byte
	VoiceInChannel  chan []byte
	VoiceOutChannel chan []byte
	InputChannel    chan []byte

	Done     bool
	lastTime time.Time
	curFPS   int

	RoomID string

	// store thing related to game
	GameMeta GameMeta

	//추가된 변수들
	audioTrack *webrtc.TrackLocalStaticSample
	videoTrack *webrtc.TrackLocalStaticSample
	peerConnection *webrtc.PeerConnection

	gateway *janus.Gateway
	session *janus.Session
	handle	*janus.Handle
}

const JANUS_URL = "ws://localhost:8188/janus"
//const JANUS_URL = "ws://192.168.56.162:8188/janus"

// NewWebRTC create
func NewWebRTC() *WebRTC {

	w := &WebRTC{
		//ID: uuid.Must(uuid.NewV4()).String(),
		ID: xid.New().String(),

		ImageChannel:    make(chan WebFrame, 30),
		AudioChannel:    make(chan []byte, 1),
		VoiceInChannel:  make(chan []byte, 1),
		VoiceOutChannel: make(chan []byte, 1),
		InputChannel:    make(chan []byte, 100),

		audioTrack: &webrtc.TrackLocalStaticSample{},
		videoTrack: &webrtc.TrackLocalStaticSample{},

	}
	return w
}

// StartClient start webrtc
func (w *WebRTC) StartClient(isMobile bool, iceCB OnIceCallback) (string, error) {
	/**
	* WebRTC pion 설정해야 하지만 일단 이 부분은 나중에 연동처리해야함.
	* 지금은 에뮬레이터를 먼저 기동 시키고 사운드와 화면을 가지고 오는지만 체크함.
	**/
	peerConnectionConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
    }
    // local webrtc peer connection
	var err error
    w.peerConnection, err = webrtc.NewPeerConnection(peerConnectionConfig)
    if err != nil {
		panic(err)
	}
	w.peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("Connection State has changed %s \n", connectionState.String())
	})

	// Create a audio track
	//audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "audio", "pion")
	w.audioTrack, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "synced-video", "synced-video")
    if err != nil {
		panic(err)
	} else if _, err = w.peerConnection.AddTrack(w.audioTrack); err != nil {
		panic(err)
	}

	// Create a video track
	//videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "video", "pion")
	w.videoTrack, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/h264"}, "synced-video", "synced-video")
    if err != nil {
		panic(err)
	} else if _, err = w.peerConnection.AddTrack(w.videoTrack); err != nil {
		panic(err)
	}


	offer, err := w.peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(w.peerConnection)

	if err = w.peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	<-gatherComplete

	log.Println(JANUS_URL)
	w.gateway, err = janus.Connect(JANUS_URL)
	if err != nil {
		panic(err)
	}

	w.session, err = w.gateway.Create()
	if err != nil {
		panic(err)
	}

	w.handle, err = w.session.Attch("janus.plugin.videoroom")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			if _, keepAliveErr := w.session.KeepAlive(); keepAliveErr != nil {
				panic(keepAliveErr)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	go watchHandle(w.handle)


	msg, err := w.handle.Message(map[string]interface{}{
		"request": "join",
		"ptype":   "publisher",
		"room":    1234,
		"display": "webrtc game",
	}, nil)
	if err != nil {
		panic(err)
	}

	feedId := msg.Plugindata.Data["id"]

	msg, err = w.handle.Message(map[string]interface{}{
		"request": "publish",
		"audio":   true,
		"video":   true,
		"data":    false,
		"bitrate": 128000,
		"bitrate_cap": true,
	}, map[string]interface{}{
		"type":    "offer",
		"sdp":     w.peerConnection.LocalDescription().SDP,
		"trickle": false,
	})
	if err != nil {
		panic(err)
	}

	if msg.Jsep != nil {
		err = w.peerConnection.SetRemoteDescription(webrtc.SessionDescription{
			Type: webrtc.SDPTypeAnswer,
			SDP:  msg.Jsep["sdp"].(string),
		})
		if err != nil {
			panic(err)
		}
		// Start pushing buffers on these tracks
		w.startStreaming(w.videoTrack, w.audioTrack)
		// 강제로 연결 처리함.
		w.isConnected = true

		//RTP forward
		_, err = w.handle.Message(map[string]interface{}{
			"request": "rtp_forward",
			"room": 1234,
			"publisher_id": feedId,
			"host": "192.168.56.168",
			"secret": "adminpwd",
			"audio_port": 5006,
			"video_port": 5011,
		}, nil)
		if err != nil {
			log.Println(err.Error())
			//panic(err)
		}

	}

	return "", nil
}


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

func (w *WebRTC) AttachRoomID(roomID string) {
	w.RoomID = roomID
}

func (w *WebRTC) SetRemoteSDP(remoteSDP string) error {
	/**
	var answer webrtc.SessionDescription
	err := Decode(remoteSDP, &answer)
	if err != nil {
		log.Println("Decode remote sdp from peer failed")
		return err
	}
	err = w.peerConnection.SetRemoteDescription(answer)
	if err != nil {
		log.Println("Set remote description from peer failed")
		return err
	}
	log.Println("Set Remote Description")
	**/
	return nil
}

func (w *WebRTC) AddCandidate(candidate string) error {
	/**
	var iceCandidate webrtc.ICECandidateInit
	err := Decode(candidate, &iceCandidate)
	if err != nil {
		log.Println("Decode Ice candidate from peer failed")
		return err
	}
	log.Println("Decoded Ice: " + iceCandidate.Candidate)
	err = w.peerConnection.AddICECandidate(iceCandidate)
	if err != nil {
		log.Println("Add Ice candidate from peer failed")
		return err
	}
	log.Println("Add Ice Candidate: " + iceCandidate.Candidate)
	**/
	return nil
}

// StopClient disconnect
func (w *WebRTC) StopClient() {
	// if stopped, bypass
	if w.isConnected == false {
		return
	}

	log.Println("===StopClient===")
	w.isConnected = false
	if w.peerConnection != nil {
		w.peerConnection.Close()
	}
	w.peerConnection = nil
	close(w.InputChannel)
	// webrtc is producer, so we close
	// NOTE: ImageChannel is waiting for input. Close in writer is not correct for this
	close(w.ImageChannel)
	close(w.AudioChannel)
	close(w.VoiceInChannel)
	close(w.VoiceOutChannel)
}

// IsConnected comment
func (w *WebRTC) IsConnected() bool {
	return w.isConnected
}


//func (w *WebRTC) startStreaming(vp8Track *webrtc.Track, opusTrack *webrtc.Track) {
func (w *WebRTC) startStreaming(vp8Track *webrtc.TrackLocalStaticSample, opusTrack *webrtc.TrackLocalStaticSample) {

	log.Println("Start streaming")
	// receive frame buffer
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered from err", r)
				log.Println(debug.Stack())
			}
		}()

		/* send
		type WebFrame struct {
			Data      []byte
			Timestamp uint32
		}
		*/
		for data := range w.ImageChannel {
			/*
			packets := vp8Track.Packetizer().Packetize(data.Data, 1)
			for _, p := range packets {
				p.Header.Timestamp = data.Timestamp
				err := vp8Track.WriteRTP(p)
				if err != nil {
					log.Println("Warn: Err write sample: ", err)
					break
				}
			}
			*/
			//log.Println("ImageChannel")
			err := vp8Track.WriteSample(media.Sample{
				Data:    data.Data,
				Duration: time.Duration(data.Timestamp),
			})
			if err != nil {
				log.Println("Warn: Err write sample: ", err)
			}

		}
	}()

	// send audio
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered from err", r)
				log.Println(debug.Stack())
			}
		}()

		for data := range w.AudioChannel {
			if !w.isConnected {
				return
			}
			audioDuration := time.Duration(20) * time.Millisecond
			err := opusTrack.WriteSample(media.Sample{
				Data:    data,
				Duration: audioDuration,

				//Samples: uint32(config.AUDIO_FRAME / config.AUDIO_CHANNELS),
			})
			if err != nil {
				log.Println("Warn: Err write sample: ", err)
			}
		}
	}()

	// send voice
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered from err", r)
				log.Println(debug.Stack())
			}
		}()

		for data := range w.VoiceOutChannel {
			if !w.isConnected {
				return
			}

			//_, err := opusTrack.Write(data)
			err := opusTrack.WriteSample(media.Sample{
				Data:    data,
			})
			if err != nil {
				log.Println("Warn: Err write sample: ", err)
			}

		}
	}()

}

func (w *WebRTC) calculateFPS() int {
	elapsedTime := time.Now().Sub(w.lastTime)
	w.lastTime = time.Now()
	curFPS := time.Second / elapsedTime
	w.curFPS = int(float32(w.curFPS)*0.9 + float32(curFPS)*0.1)
	return w.curFPS
}
