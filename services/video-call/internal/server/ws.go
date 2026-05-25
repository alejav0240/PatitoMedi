package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"patitomedi/video-call/internal/hub"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade: %v", err)
		return
	}

	peer := &hub.Peer{Send: make(chan []byte, 64)}
	var roomID string

	defer func() {
		conn.Close()
		if roomID != "" && peer.UserID != "" {
			s.hub.Leave(roomID, peer.UserID)
		}
		close(peer.Send)
	}()

	go func() {
		for data := range peer.Send {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		}
	}()

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg hub.Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			sendError(peer, "invalid json")
			continue
		}

		switch msg.Type {
		case "join-room":
			if msg.AppointmentID == "" || msg.UserID == "" {
				sendError(peer, "appointmentId and userId required")
				continue
			}
			roomID = msg.AppointmentID
			peer.UserID = msg.UserID
			s.hub.Join(roomID, peer)
			activeRoomsGauge.Set(float64(s.hub.ActiveRooms()))
			activeParticipantsGauge.Set(float64(s.hub.ActiveParticipants()))

		case "leave-room":
			if roomID != "" {
				s.hub.Leave(roomID, peer.UserID)
				activeRoomsGauge.Set(float64(s.hub.ActiveRooms()))
				activeParticipantsGauge.Set(float64(s.hub.ActiveParticipants()))
				roomID = ""
			}

		case "offer", "answer", "ice-candidate":
			if roomID == "" {
				sendError(peer, "not in a room")
				continue
			}
			s.hub.Route(roomID, &msg, raw, peer.UserID)

		case "call-ended":
			if roomID != "" {
				s.hub.Leave(roomID, peer.UserID)
				activeRoomsGauge.Set(float64(s.hub.ActiveRooms()))
				activeParticipantsGauge.Set(float64(s.hub.ActiveParticipants()))
				roomID = ""
			}

		default:
			sendError(peer, "unknown message type: "+msg.Type)
		}
	}
}

func sendError(p *hub.Peer, msg string) {
	data, _ := json.Marshal(map[string]string{"type": "error", "message": msg})
	select {
	case p.Send <- data:
	default:
	}
}
