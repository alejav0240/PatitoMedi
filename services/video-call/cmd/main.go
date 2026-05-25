package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"patitomedi/video-call/internal/events"
	"patitomedi/video-call/internal/hub"
	"patitomedi/video-call/internal/server"
	"patitomedi/video-call/internal/store"
)

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	port := getenv("PORT", "8080")
	redisAddr := getenv("REDIS_ADDR", "localhost:6379")
	kafkaBrokers := strings.Split(getenv("KAFKA_BROKERS", "localhost:9092"), ",")

	rs := store.New(redisAddr)
	if err := rs.Ping(context.Background()); err != nil {
		log.Printf("redis ping failed: %v (continuing)", err)
	}

	producer := events.New(kafkaBrokers)
	defer producer.Close()

	h := hub.New()

	h.OnRoomCreated = func(appointmentID string) {
		rs.RoomOpened(context.Background(), appointmentID)
	}

	h.OnPeerJoined = func(appointmentID, userID string, peerCount int) {
		rs.PeerJoined(context.Background(), appointmentID, userID)
		if peerCount == 2 {
			server.CallsTotalCounter.Inc()
			producer.CallStarted(context.Background(), appointmentID, []string{userID})
		}
	}

	h.OnPeerLeft = func(appointmentID, userID string, peerCount int) {
		rs.PeerLeft(context.Background(), appointmentID, userID)
		if peerCount == 0 {
			producer.CallEnded(context.Background(), appointmentID, userID)
		}
	}

	h.OnRoomClosed = func(appointmentID string) {
		rs.RoomClosed(context.Background(), appointmentID)
	}

	srv := server.New(h)
	log.Printf("video-call service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		log.Fatal(err)
	}
}
