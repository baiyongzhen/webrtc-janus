package worker

import (

	"log"
	"os"
	"path"

	"example.com/webrtc-game/pkg/config"
	"example.com/webrtc-game/pkg/util"
	storage "example.com/webrtc-game/pkg/worker/cloud-storage"
	"example.com/webrtc-game/pkg/worker/room"
)

type Handler struct {
	cfg             config.Config
	// Rooms map : RoomID -> Room
	rooms map[string]*room.Room
	// ID of the current server globalwise
	serverID string
	// onlineStorage is client accessing to online storage (GCP)
	onlineStorage *storage.Client
	// sessions handles all sessions server is handler (key is sessionID)
	sessions map[string]*Session

    //임시적으로 input keymap 테스트
	roomID string
}

// NewHandler returns a new server
func NewHandler(cfg config.Config) *Handler {
	// Create offline storage folder
	createOfflineStorage()
	// Init online storage
	onlineStorage := storage.NewInitClient()
	return &Handler{
		rooms:           map[string]*room.Room{},
		sessions:        map[string]*Session{},
		cfg:             cfg,
		onlineStorage:   onlineStorage,
	}
}

//임시적으로 input keymap 테스트
func (h *Handler) InputKeyboard() {
	roomID := h.roomID
	log.Println("room id:", roomID)
	h.rooms[roomID].InputKeyboard()
}

// Run starts a Handler running logic
func (h *Handler) Run() {
	h.RouteCoordinator()
}

func (h *Handler) Close() {
	// Close all room
}


// getRoom returns room from roomID
func (h *Handler) getRoom(roomID string) *room.Room {
	room, ok := h.rooms[roomID]
	if !ok {
		return nil
	}

	return room
}

// createNewRoom creates a new room
// Return nil in case of room is existed
func (h *Handler) createNewRoom(gameName string, roomID string, videoEncoderType string) *room.Room {
	// If the roomID is empty,
	// or the roomID doesn't have any running sessions (room was closed)
	// we spawn a new room
	if roomID == "" || !h.isRoomRunning(roomID) {
		room := room.NewRoom(roomID, gameName, videoEncoderType, h.onlineStorage, h.cfg)
		// TODO: Might have race condition
		h.rooms[room.ID] = room

		//임시적으로 input keymap 테스트
		h.roomID = room.ID

		return room
	}

	return nil
}

// isRoomRunning check if there is any running sessions.
// TODO: If we remove sessions from room anytime a session is closed, we can check if the sessions list is empty or not.
func (h *Handler) isRoomRunning(roomID string) bool {
	// If no roomID is registered
	room, ok := h.rooms[roomID]
	if !ok {
		return false
	}

	return room.IsRunningSessions()
}


func createOfflineStorage() {
	dir, _ := path.Split(util.GetSavePath("dummy"))
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Println("Failed to create offline storage, err: ", err)
	}
}
