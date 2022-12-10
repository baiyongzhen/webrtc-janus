package room

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"example.com/webrtc-game/pkg/config"
	"example.com/webrtc-game/pkg/emulator"
	"example.com/webrtc-game/pkg/emulator/libretro/nanoarch"
	"example.com/webrtc-game/pkg/util"
	"example.com/webrtc-game/pkg/util/gamelist"
	"example.com/webrtc-game/pkg/webrtc"
	storage "example.com/webrtc-game/pkg/worker/cloud-storage"
)

type Room struct {
	ID string

	// imageChannel is image stream received from director
	imageChannel <-chan nanoarch.GameFrame
	// audioChannel is audio stream received from director
	audioChannel <-chan []int16
	// inputChannel is input stream send to director. This inputChannel is combined
	// input from webRTC + connection info (player indexc)
	inputChannel chan<- nanoarch.InputEvent
	// voiceInChannel is voice stream received from users
	voiceInChannel chan []byte
	// voiceOutChannel is voice stream broadcasted to all users
	voiceOutChannel chan []byte
	voiceSample     [][]byte
	// State of room
	IsRunning bool
	// Done channel is to fire exit event when room is closed
	Done chan struct{}
	// List of peerconnections in the room
	rtcSessions []*webrtc.WebRTC
	// NOTE: Not in use, lock rtcSessions
	sessionsLock *sync.Mutex
	// Director is emulator
	director emulator.CloudEmulator
	// Cloud storage to store room state online
	onlineStorage *storage.Client
	// GameName
	gameName string
	// Meta of game
	//meta emulator.Meta

}

const separator = "___"

// TODO: Remove after fully migrate
const oldSeparator = "|"

const SocketAddrTmpl = "/tmp/cloudretro-retro-%s.sock"
const bufSize = 245969

// NewVideoImporter return image Channel from stream
func NewVideoImporter(roomID string) chan nanoarch.GameFrame {
	sockAddr := fmt.Sprintf(SocketAddrTmpl, roomID)
	imgChan := make(chan nanoarch.GameFrame)

	l, err := net.Listen("unix", sockAddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	log.Println("Creating uds server", sockAddr)
	go func(l net.Listener) {
		defer l.Close()

		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}
		defer conn.Close()

		log.Println("Received new conn")
		log.Println("Spawn Importer")

		fullbuf := make([]byte, bufSize*2)
		fullbuf = fullbuf[:0]

		for {
			// TODO: Not reallocate
			buf := make([]byte, bufSize)
			l, err := conn.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("error: %v", err)
				}
				continue
			}

			buf = buf[:l]
			fullbuf = append(fullbuf, buf...)
			if len(fullbuf) >= bufSize {
				bufs := bytes.NewBuffer(fullbuf)
				dec := gob.NewDecoder(bufs)

				frame := nanoarch.GameFrame{}
				err := dec.Decode(&frame)
				if err != nil {
					log.Fatalf("%v", err)
				}
				imgChan <- frame
				fullbuf = fullbuf[bufSize:len(fullbuf)]
			}
		}
	}(l)

	return imgChan
}

// NewRoom creates a new room
func NewRoom(roomID string, gameName string,
	videoEncoderType string,
	onlineStorage *storage.Client, cfg config.Config) *Room {
	// If no roomID is given, generate it from gameName
	// If the is roomID, get gameName from roomID
	if roomID == "" {
		roomID = generateRoomID(gameName)
	} else {
		gameName = getGameNameFromRoomID(roomID)
		log.Println("Get Gamename from RoomID", gameName)
	}
	gameInfo := gamelist.GetGameInfoFromName(gameName)

	log.Println("roomID: ", roomID)
	log.Println("gameName: ", gameName)
	log.Println("gameInfo: ", gameInfo)


	inputChannel := make(chan nanoarch.InputEvent, 100)

	room := &Room{
		ID: roomID,

		inputChannel:    inputChannel,
		imageChannel:    nil,
		voiceInChannel:  make(chan []byte, 1),
		voiceOutChannel: make(chan []byte, 1),
		rtcSessions:     []*webrtc.WebRTC{},
		sessionsLock:    &sync.Mutex{},
		IsRunning:       true,
		onlineStorage:   onlineStorage,

		Done: make(chan struct{}, 1),
	}

	// Check if room is on local storage, if not, pull from GCS to local storage
	go func(game gamelist.GameInfo, roomID string) {
		// Check room is on local or fetch from server
		savepath := util.GetSavePath(roomID)
		log.Println("Check ", savepath, " on online storage : ", room.isGameOnLocal(savepath))
		if err := room.saveOnlineRoomToLocal(roomID, savepath); err != nil {
			log.Printf("Warn: Room %s is not in online storage, error %s", roomID, err)
		}

		// If not then load room or create room from local.
		log.Printf("Room %s started. GamePath: %s, GameName: %s, WithGame: %t", roomID, game.Path, game.Name, cfg.WithoutGame)

		// Spawn new emulator based on gameName and plug-in all channels
		emuName, _ := config.FileTypeToEmulator[game.Type]

		if cfg.WithoutGame {
			// Run without game, image stream is communicated over unixsocket
			imageChannel := NewVideoImporter(roomID)
			director, _, audioChannel := nanoarch.Init(emuName, roomID, false, inputChannel)
			room.imageChannel = imageChannel
			room.director = director
			room.audioChannel = audioChannel
		} else {
			// Run without game, image stream is communicated over image channel
			director, imageChannel, audioChannel := nanoarch.Init(emuName, roomID, true, inputChannel)
			room.imageChannel = imageChannel
			room.director = director
			room.audioChannel = audioChannel
		}

		gameMeta := room.director.LoadMeta(game.Path)

		// nwidth, nheight are the webRTC output size.
		// There are currently two approach
		var nwidth, nheight int
		if cfg.EnableAspectRatio {
			baseAspectRatio := float64(gameMeta.BaseWidth) / float64(gameMeta.Height)
			nwidth, nheight = resizeToAspect(baseAspectRatio, cfg.Width, cfg.Height)
			log.Printf("Viewport size will be changed from %dx%d (%f) -> %dx%d", cfg.Width, cfg.Height,
				baseAspectRatio, nwidth, nheight)
		} else {
			nwidth, nheight = gameMeta.BaseWidth, gameMeta.BaseHeight
			log.Printf("Viewport custom size is disabled, base size will be used instead %dx%d", nwidth, nheight)
		}
		if cfg.Scale > 1 {
			nwidth, nheight = nwidth*cfg.Scale, nheight*cfg.Scale
			log.Printf("Viewport size has scaled to %dx%d", nwidth, nheight)
		}

		log.Println("meta: ", gameMeta)

		// set resulting game frame size considering
		// its orientation
		encoderW, encoderH := nwidth, nheight
		if gameMeta.Rotation.IsEven {
			encoderW, encoderH = nheight, nwidth
		}

		room.director.SetViewport(encoderW, encoderH)

		// Spawn video and audio encoding for webRTC
		go room.startVideo(encoderW, encoderH, videoEncoderType)
		go room.startAudio(gameMeta.AudioSampleRate)
		go room.startVoice()
		room.director.Start()

		log.Printf("Room %s ended", roomID)

		// TODO: do we need GC, we can remove it
		runtime.GC()
	}(gameInfo, roomID)


	return room
}

func resizeToAspect(ratio float64, sw int, sh int) (dw int, dh int) {
	// ratio is always > 0
	dw = int(math.Round(float64(sh)*ratio/2) * 2)
	dh = sh
	if dw > sw {
		dw = sw
		dh = int(math.Round(float64(sw)/ratio/2) * 2)
	}
	return
}

// getGameNameFromRoomID parse roomID to get roomID and gameName
func getGameNameFromRoomID(roomID string) string {
	parts := strings.Split(roomID, separator)
	if len(parts) > 1 {
		return parts[1]
	}
	// TODO: Remove when fully migrate
	parts = strings.Split(roomID, oldSeparator)
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// generateRoomID generate a unique room ID containing 16 digits
func generateRoomID(gameName string) string {
	// RoomID contains random number + gameName
	// Next time when we only get roomID, we can launch game based on gameName
	roomID := strconv.FormatInt(rand.Int63(), 16) + separator + gameName
	return roomID
}

func (r *Room) isGameOnLocal(savepath string) bool {
	_, err := os.Open(savepath)
	return err == nil
}

func (r *Room) AddConnectionToRoom(peerconnection *webrtc.WebRTC) {
	peerconnection.AttachRoomID(r.ID)
	r.rtcSessions = append(r.rtcSessions, peerconnection)

	//임시적으로 처리함.
	//go r.startWebRTCSession(peerconnection)
}

//임시적으로 input keymap 테스트
func (r *Room) InputKeyboard() {
	//id := room.rtcSessions[0]
	//room.inputChannel <- nanoarch.InputEvent{RawState: []byte{0xFF, 0xFF}, ConnID: w.ID}
	//var state  []byte
	//state[0] = 64
	//state[1] = 0

	state := []byte{0x40, 0x00}
	//state := []byte{0x00, 0x40}
	for _, s := range r.rtcSessions {
		//if s.IsConnected() {
		r.inputChannel <- nanoarch.InputEvent{RawState: state, PlayerIdx: s.GameMeta.PlayerIndex, ConnID: s.ID}
		//}
	}
	time.Sleep(50 * time.Millisecond)
	state = []byte{0x00, 0x00}
	for _, s := range r.rtcSessions {
		//if s.IsConnected() {
		r.inputChannel <- nanoarch.InputEvent{RawState: state, PlayerIdx: s.GameMeta.PlayerIndex, ConnID: s.ID}
		//}
	}

}

func (r *Room) startWebRTCSession(peerconnection *webrtc.WebRTC) {
	//client voice && datachannel 정보를 받는다..
	/**
	defer func() {
		if r := recover(); r != nil {
			log.Println("Warn: Recovered when sent to close inputChannel")
		}
	}()

	log.Println("Start WebRTC session")
	go func() {

		// set up voice input and output. A room has multiple voice input and only one combined voice output.
		for voiceInput := range peerconnection.VoiceInChannel {
			// NOTE: when room is no longer running. InputChannel needs to have extra event to go inside the loop
			if peerconnection.Done || !peerconnection.IsConnected() || !r.IsRunning {
				break
			}

			if peerconnection.IsConnected() {
				r.voiceInChannel <- voiceInput
			}

		}
	}()

	// bug: when inputchannel here = nil , skip and finish
	for input := range peerconnection.InputChannel {
		// NOTE: when room is no longer running. InputChannel needs to have extra event to go inside the loop
		if peerconnection.Done || !peerconnection.IsConnected() || !r.IsRunning {
			break
		}

		if peerconnection.IsConnected() {
			select {
			case r.inputChannel <- nanoarch.InputEvent{RawState: input, PlayerIdx: peerconnection.GameMeta.PlayerIndex, ConnID: peerconnection.ID}:
			default:
			}
		}
	}

	log.Println("Peerconn done")


	**/


}


// saveOnlineRoomToLocal save online room to local
func (r *Room) saveOnlineRoomToLocal(roomID string, savepath string) error {
	log.Println("Check if game is on cloud storage")
	// If the game is not on local server
	// Try to load from gcloud
	data, err := r.onlineStorage.LoadFile(roomID)
	if err != nil {
		return err
	}
	// Save the data fetched from gcloud to local server
	ioutil.WriteFile(savepath, data, 0644)

	return nil
}

func (r *Room) IsRunningSessions() bool {
	// If there is running session
	for _, s := range r.rtcSessions {
		if s.IsConnected() {
			return true
		}
	}

	return false
}
