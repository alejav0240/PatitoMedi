package hub

import (
	"encoding/json"
	"sync"
)

// Message is the generic WebSocket envelope.
type Message struct {
	Type          string          `json:"type"`
	AppointmentID string          `json:"appointmentId,omitempty"`
	UserID        string          `json:"userId,omitempty"`
	To            string          `json:"to,omitempty"`
	SDP           string          `json:"sdp,omitempty"`
	Candidate     string          `json:"candidate,omitempty"`
	SdpMid        string          `json:"sdpMid,omitempty"`
	SdpMLineIndex *int            `json:"sdpMLineIndex,omitempty"`
	Message       string          `json:"message,omitempty"`
	Raw           json.RawMessage `json:"-"`
}

// Peer represents a connected WebSocket client.
type Peer struct {
	UserID string
	Send   chan []byte
}

// Room holds the peers for one appointment.
type Room struct {
	AppointmentID string
	peers         map[string]*Peer // keyed by userID
	mu            sync.RWMutex
}

func newRoom(appointmentID string) *Room {
	return &Room{AppointmentID: appointmentID, peers: make(map[string]*Peer)}
}

func (r *Room) add(p *Peer) {
	r.mu.Lock()
	r.peers[p.UserID] = p
	r.mu.Unlock()
}

func (r *Room) remove(userID string) {
	r.mu.Lock()
	delete(r.peers, userID)
	r.mu.Unlock()
}

func (r *Room) size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.peers)
}

func (r *Room) peer(userID string) (*Peer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.peers[userID]
	return p, ok
}

// broadcast sends data to all peers except the sender.
func (r *Room) broadcast(data []byte, senderID string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for id, p := range r.peers {
		if id != senderID {
			select {
			case p.Send <- data:
			default:
			}
		}
	}
}

// Hub manages all rooms.
type Hub struct {
	rooms map[string]*Room
	mu    sync.RWMutex

	// Callbacks for external integrations (Redis, Kafka).
	OnRoomCreated func(appointmentID string)
	OnRoomClosed  func(appointmentID string)
	OnPeerJoined  func(appointmentID, userID string, peerCount int)
	OnPeerLeft    func(appointmentID, userID string, peerCount int)
}

func New() *Hub {
	return &Hub{rooms: make(map[string]*Room)}
}

func (h *Hub) getOrCreate(appointmentID string) (*Room, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	r, ok := h.rooms[appointmentID]
	if !ok {
		r = newRoom(appointmentID)
		h.rooms[appointmentID] = r
		return r, true // created
	}
	return r, false
}

func (h *Hub) Join(appointmentID string, p *Peer) {
	room, created := h.getOrCreate(appointmentID)
	room.add(p)
	count := room.size()
	if created && h.OnRoomCreated != nil {
		h.OnRoomCreated(appointmentID)
	}
	if h.OnPeerJoined != nil {
		h.OnPeerJoined(appointmentID, p.UserID, count)
	}
}

func (h *Hub) Leave(appointmentID, userID string) {
	h.mu.Lock()
	room, ok := h.rooms[appointmentID]
	h.mu.Unlock()
	if !ok {
		return
	}
	room.remove(userID)
	count := room.size()
	if h.OnPeerLeft != nil {
		h.OnPeerLeft(appointmentID, userID, count)
	}
	if count == 0 {
		h.mu.Lock()
		delete(h.rooms, appointmentID)
		h.mu.Unlock()
		if h.OnRoomClosed != nil {
			h.OnRoomClosed(appointmentID)
		}
	}
}

// Route delivers a message to a specific peer or broadcasts it.
func (h *Hub) Route(appointmentID string, msg *Message, data []byte, senderID string) {
	h.mu.RLock()
	room, ok := h.rooms[appointmentID]
	h.mu.RUnlock()
	if !ok {
		return
	}
	if msg.To != "" {
		if p, found := room.peer(msg.To); found {
			select {
			case p.Send <- data:
			default:
			}
		}
		return
	}
	room.broadcast(data, senderID)
}

// ActiveRooms returns the current room count.
func (h *Hub) ActiveRooms() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms)
}

// ActiveParticipants returns total connected peers across all rooms.
func (h *Hub) ActiveParticipants() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	total := 0
	for _, r := range h.rooms {
		total += r.size()
	}
	return total
}
