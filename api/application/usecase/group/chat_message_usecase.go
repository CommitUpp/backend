package group

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/CommitUpp/backend/api/domain/repository"
)

const maxChatMessageLength = 1000

var (
	ErrChatMessageRoomIDRequired  = errors.New("chat message room ID is required")
	ErrChatMessageContentRequired = errors.New("chat message content is required")
	ErrChatMessageContentTooLong  = errors.New("chat message content is too long")
	ErrChatRoomNotFound           = repository.ErrChatRoomNotFound
	ErrNotGroupMember             = repository.ErrNotGroupMember
)

type ChatMessageUsecase interface {
	CreateChatMessage(ctx context.Context, input CreateChatMessageInput) (repository.ChatMessage, error)
}

type CreateChatMessageInput struct {
	RoomID  string
	UserID  string
	Content string
}

type chatMessageUsecase struct {
	chatMessageRepo repository.ChatMessageRepository
}

func NewChatMessageUsecase(chatMessageRepo repository.ChatMessageRepository) ChatMessageUsecase {
	return &chatMessageUsecase{chatMessageRepo: chatMessageRepo}
}

func (u *chatMessageUsecase) CreateChatMessage(
	ctx context.Context,
	input CreateChatMessageInput,
) (repository.ChatMessage, error) {
	if input.UserID == "" {
		return repository.ChatMessage{}, ErrUserIDRequired
	}
	if input.RoomID == "" {
		return repository.ChatMessage{}, ErrChatMessageRoomIDRequired
	}

	content := strings.TrimSpace(input.Content)
	if content == "" {
		return repository.ChatMessage{}, ErrChatMessageContentRequired
	}
	if utf8.RuneCountInString(content) > maxChatMessageLength {
		return repository.ChatMessage{}, ErrChatMessageContentTooLong
	}

	return u.chatMessageRepo.CreateChatMessage(ctx, repository.CreateChatMessageInput{
		RoomID:  input.RoomID,
		UserID:  input.UserID,
		Content: content,
	})
}
