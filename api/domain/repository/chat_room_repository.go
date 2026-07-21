package repository

import (
	"context"
	"time"
)

type ChatRoomRepository interface {
	GetChatRooms(ctx context.Context, groupID string) ([]ChatRoom, error)
}

type ChatRoom struct {
	ID         string
	GroupID    string
	MovieID    *string
	MovieTitle string
	CreatedBy  string
	CreatedAt  time.Time
}
