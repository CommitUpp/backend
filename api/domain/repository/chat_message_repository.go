package repository

import (
	"context"
	"errors"
	"time"
)

var (
	ErrChatRoomNotFound = errors.New("chat room not found")
	ErrNotGroupMember   = errors.New("not an active group member")
)

type ChatMessageRepository interface {
	CreateChatMessage(ctx context.Context, input CreateChatMessageInput) (ChatMessage, error)
}

type CreateChatMessageInput struct {
	RoomID  string
	UserID  string
	Content string
}

type ChatMessage struct {
	ID        string
	RoomID    string
	UserID    string
	Content   string
	CreatedAt time.Time
}
