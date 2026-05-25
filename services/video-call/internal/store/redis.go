package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const roomTTL = 2 * time.Hour

type RedisStore struct {
	client *redis.Client
}

func New(addr string) *RedisStore {
	return &RedisStore{
		client: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (s *RedisStore) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

func roomKey(appointmentID string) string {
	return fmt.Sprintf("room:%s", appointmentID)
}

func participantsKey(appointmentID string) string {
	return fmt.Sprintf("room:%s:participants", appointmentID)
}

// RoomOpened sets room status and resets TTL.
func (s *RedisStore) RoomOpened(ctx context.Context, appointmentID string) error {
	pipe := s.client.Pipeline()
	pipe.HSet(ctx, roomKey(appointmentID), "status", "active", "started_at", time.Now().UTC().Format(time.RFC3339))
	pipe.Expire(ctx, roomKey(appointmentID), roomTTL)
	_, err := pipe.Exec(ctx)
	return err
}

// PeerJoined adds a participant and refreshes TTL.
func (s *RedisStore) PeerJoined(ctx context.Context, appointmentID, userID string) error {
	pipe := s.client.Pipeline()
	pipe.SAdd(ctx, participantsKey(appointmentID), userID)
	pipe.Expire(ctx, participantsKey(appointmentID), roomTTL)
	_, err := pipe.Exec(ctx)
	return err
}

// PeerLeft removes a participant.
func (s *RedisStore) PeerLeft(ctx context.Context, appointmentID, userID string) error {
	return s.client.SRem(ctx, participantsKey(appointmentID), userID).Err()
}

// RoomClosed deletes all room keys.
func (s *RedisStore) RoomClosed(ctx context.Context, appointmentID string) error {
	return s.client.Del(ctx, roomKey(appointmentID), participantsKey(appointmentID)).Err()
}
