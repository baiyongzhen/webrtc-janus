package worker

import (
	"log"

	"example.com/webrtc-game/pkg/util"
	"example.com/webrtc-game/pkg/webrtc"
	"example.com/webrtc-game/pkg/worker/room"
)

// RouteCoordinator are all routes server received from coordinator
// Socket 정보로 전달함.
func (h *Handler) RouteCoordinator() {
	// co
	// 00. serverID
	// 01. initwebrtc (webrtc,janus 연동이 필요함.)
	peerconnection := webrtc.NewWebRTC()
	var initPacket struct {
		IsMobile bool `json:"is_mobile"`
	}
	initPacket.IsMobile = false

	//localSession, err := peerconnection.StartClient(
	_, err := peerconnection.StartClient(
		initPacket.IsMobile,
		func(candidate string) {
			// send back candidate string to browser
			//oClient.Send(cws.WSPacket{
			//	ID:        "candidate",
			//	Data:      candidate,
			//	SessionID: resp.SessionID,
			//}, nil)
		},
	)

	// Create new sessions when we have new peerconnection initialized
	session := &Session{
		peerconnection: peerconnection,
	}

	//h.sessions[resp.SessionID] = session
	sessionID := "aaaaaaaaaa"
	h.sessions[sessionID] = session
	log.Println("Start peerconnection", sessionID)

	if err != nil {
		log.Println("Error: Cannot create new webrtc session", err)
		//return cws.EmptyPacket
	}

	// 02. answer
	// 03. candidate
	// 04. start

	/**
	//임의적으로 처리한 부분임.
	01. initwebrtc 에서 클라이언트로 전달해 준 값을 받아야 하지만
	log.Println("Received a start request from coordinator")
	session := h.getSession(resp.SessionID)
	if session == nil {
		log.Printf("Error: No session for ID: %s\n", resp.SessionID)
		return cws.EmptyPacket
	}
	**/

	//session := h.getSession(sessionID)
	peerconnection = session.peerconnection

	var startPacket struct {
		GameName string `json:"game_name"`
		IsMobile bool   `json:"is_mobile"`
	}
	//startPacket.GameName = "Contra (U)"
	//startPacket.GameName = "Street Blaster II Pro"
	startPacket.GameName = "Street Fighter 3"
	startPacket.IsMobile = false
	//roomID := "Contra"
	roomID := ""
	playerIndex := 0

	log.Println(startPacket)
	//room := h.startGameHandler(startPacket.GameName, resp.RoomID, resp.PlayerIndex, peerconnection, util.GetVideoEncoder(startPacket.IsMobile))
	room := h.startGameHandler(startPacket.GameName, roomID, playerIndex, peerconnection,util.GetVideoEncoder(startPacket.IsMobile))

	// TODO: can data race
	h.rooms[room.ID] = room

	// 05. quit
	// 06. save
	// 07. load
	// 08. playerIdx
	// 09. multitap
	// 10. terminateSession

}

// startGameHandler starts a game if roomID is given, if not create new room
//func (h *Handler) startGameHandler(gameName, existedRoomID string, playerIndex int, videoEncoderType string) *room.Room {
func (h *Handler) startGameHandler(gameName, existedRoomID string, playerIndex int, peerconnection *webrtc.WebRTC, videoEncoderType string) *room.Room {

	log.Println("Starting gameName:", gameName)
	log.Println("Starting RoomID:", existedRoomID)
	log.Println("Starting videoEncoderType:", videoEncoderType)

	room := h.getRoom(existedRoomID)
	// If room is not running
	if room == nil {
		//log.Println("Got Room from local ", room, " ID: ", existedRoomID)
		// Create new room and update player index
		room = h.createNewRoom(gameName, existedRoomID, videoEncoderType)
		// Wait for done signal from room
	}

	// Attach peerconnection to room. If PC is already in room, don't detach
	// Attach peerconnection to room. If PC is already in room, don't detach
	/*
	log.Println("Is PC in room", room.IsPCInRoom(peerconnection))
	if !room.IsPCInRoom(peerconnection) {
		h.detachPeerConn(peerconnection)
		room.AddConnectionToRoom(peerconnection)
	}
	*/
	room.AddConnectionToRoom(peerconnection)

	// Register room to coordinator if we are connecting to coordinator

	return room
}
