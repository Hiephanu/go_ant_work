package server

import (
	"encoding/json"
	"fmt"
	"go_ant_work/internal/structs"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Presence struct {
	UserID string
	Host   bool
	Conn   *websocket.Conn
}

type Room struct {
	Id        string
	Name      string
	Password  string
	Presences []Presence
	HostID    string
}

type RoomTracker struct {
	Mutex sync.RWMutex
	Rooms map[string]*Room
}

// Init initializes the RoomTracker by creating the map structure.
func NewTrackerRoom() *RoomTracker {
	return &RoomTracker{
		Mutex: sync.RWMutex{},
		Rooms: make(map[string]*Room),
	}
}

func (r *RoomTracker) CreateRoom(name string, password string, userId string) structs.CreateRoomResponse {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	roomId := uuid.NewString()
	room := &Room{
		Id:       roomId,
		Name:     name,
		Password: password,
	}
	r.Rooms[roomId] = room

	presence := Presence{
		UserID: userId,
		Host:   true,
		Conn:   nil,
	}

	r.Rooms[roomId].Presences = []Presence{
		presence,
	}

	fmt.Println("Room", r.Rooms[roomId])
	return structs.CreateRoomResponse{
		RoomID:   roomId,
		RoomName: name,
	}
}

func (rt *RoomTracker) AddConnToUser(userID string, roomID string, conn *websocket.Conn) error {
	room, exists := rt.Rooms[roomID]
	if !exists {
		return fmt.Errorf("room %s does not exist", roomID)
	}

	// Initialize Presences if nil
	if room.Presences == nil {
		room.Presences = make([]Presence, 0) // Adjust to your actual type
	}

	for _, presence := range room.Presences {
		if presence.UserID == userID {
			fmt.Println("add to userId")
			presence.Conn = conn
			return nil
		}
	}

	fmt.Println("Test add conn ", roomID, conn)
	return fmt.Errorf("Presence not found")
}

func (rt *RoomTracker) SendToRoom(roomId string, envelop Envelop) {
	fmt.Println("Send to other socket")
	for _, p := range rt.Rooms[roomId].Presences {
		if p.Conn != nil {
			msg, err := json.Marshal(envelop)

			if err != nil {
				fmt.Println("Error to marshal envelop", err)
			}

			err = p.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Println("Error to send envelop", err)
			}
		}
	}
}

// GetPresences returns the list of Presences in a room.
// If the room does not exist, it returns an empty slice.
func (r *RoomTracker) GetPresences(roomID string) []Presence {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	presences := r.Rooms[roomID].Presences

	return presences
}

// AddPresence adds a new Presence to the room.
func (r *RoomTracker) JoinRoom(roomID string, useId string, password string, offer interface{}) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if r.Rooms[roomID].Password != password {
		return fmt.Errorf("Password wrong in room")
	}

	presence := &Presence{
		UserID: useId,
		Conn:   nil,
	}

	r.Rooms[roomID].Presences = append(r.Rooms[roomID].Presences, *presence)

	if len(r.Rooms[roomID].Presences) > 1 {
		offerEnvelop := &Envelop{
			Type:    Offer,
			Payload: offer,
		}
		r.SendToRoom(roomID, *offerEnvelop)
	}

	envelop := &Envelop{
		Type:    JoinRoom,
		Payload: presence,
	}

	if len(r.Rooms[roomID].Presences) > 1 {
		r.SendToRoom(roomID, *envelop)
	}

	return nil
}

func (r *RoomTracker) RemovePresence(roomID string, hostID string, userID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	room, roomExists := r.Rooms[roomID]
	if !roomExists {
		return
	}

	if room.HostID != hostID {
		return
	}

	presences := room.Presences
	for i, p := range presences {
		if p.UserID == userID {
			room.Presences = append(presences[:i], presences[i+1:]...)
			return
		}
	}

	envelop := Envelop{
		Type: KickRoom,
		Payload: structs.KickRoomEvent{
			RoomID: roomID,
			UserID: userID,
		},
	}
	r.SendToRoom(roomID, envelop)
}

func (r *RoomTracker) LeftRoom(roomID string, userID string) {
	r.Mutex.Lock()
	defer r.Mutex.Lock()

	room, roomExists := r.Rooms[roomID]
	if !roomExists {
		return
	}

	presences := room.Presences
	for i, p := range presences {
		if p.UserID == userID {
			room.Presences = append(presences[:i], presences[i+1:]...)
			break
		}
	}

	leftEnvelop := &Envelop{
		Type:    LeftRoom,
		Payload: userID,
	}

	r.SendToRoom(roomID, *leftEnvelop)
}

func (r *RoomTracker) GetPresenceByUserID(userId string, roomId string) *Presence {
	presences := r.GetPresences(roomId)

	for _, presence := range presences {
		if presence.UserID == userId {
			return &presence
		}
	}

	return nil
}

func (r *RoomTracker) IsRoomExist(roomId string) bool {
	_, roomExists := r.Rooms[roomId]
	if !roomExists {
		return false
	} else {
		return true
	}
}
