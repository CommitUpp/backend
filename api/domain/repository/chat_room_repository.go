package repository

import (
	"context"
	"errors"
	"time"
)

var (
	ErrMovieNotFound         = errors.New("movie not found")
	ErrMovieNotWatched       = errors.New("movie not watched")
	ErrChatRoomAlreadyExists = errors.New("chat room already exists")
)

type ChatRoomRepository interface {
	GetChatRooms(ctx context.Context, groupID string) ([]ChatRoom, error)
	CreateChatRoom(ctx context.Context, input CreateChatRoomInput) (ChatRoom, error)
}

type CreateChatRoomInput struct {
	GroupID   string
	MovieID   string
	CreatedBy string
}

type ChatRoom struct {
	ID         string
	GroupID    string
	MovieID    *string
	MovieTitle string
	CreatedBy  string
	CreatedAt  time.Time
}
