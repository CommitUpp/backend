package group

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/CommitUpp/backend/api/domain/repository"
)

type stubChatMessageRepository struct {
	createdInput repository.CreateChatMessageInput
	message      repository.ChatMessage
	err          error
}

func (s *stubChatMessageRepository) CreateChatMessage(
	_ context.Context,
	input repository.CreateChatMessageInput,
) (repository.ChatMessage, error) {
	s.createdInput = input
	return s.message, s.err
}

func TestChatMessageUsecaseCreateChatMessageValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input CreateChatMessageInput
		want  error
	}{
		{
			name:  "user ID is required",
			input: CreateChatMessageInput{RoomID: "room-id", Content: "message"},
			want:  ErrUserIDRequired,
		},
		{
			name:  "room ID is required",
			input: CreateChatMessageInput{UserID: "user-id", Content: "message"},
			want:  ErrChatMessageRoomIDRequired,
		},
		{
			name:  "content is required",
			input: CreateChatMessageInput{UserID: "user-id", RoomID: "room-id", Content: " \n\t "},
			want:  ErrChatMessageContentRequired,
		},
		{
			name: "content must be at most 1000 characters",
			input: CreateChatMessageInput{
				UserID:  "user-id",
				RoomID:  "room-id",
				Content: strings.Repeat("あ", maxChatMessageLength+1),
			},
			want: ErrChatMessageContentTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &stubChatMessageRepository{}
			usecase := NewChatMessageUsecase(repo)

			_, err := usecase.CreateChatMessage(context.Background(), tt.input)
			if !errors.Is(err, tt.want) {
				t.Fatalf("CreateChatMessage() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestChatMessageUsecaseCreateChatMessageTrimsAndPersistsContent(t *testing.T) {
	t.Parallel()

	want := repository.ChatMessage{ID: "message-id"}
	repo := &stubChatMessageRepository{message: want}
	usecase := NewChatMessageUsecase(repo)

	got, err := usecase.CreateChatMessage(context.Background(), CreateChatMessageInput{
		UserID:  "user-id",
		RoomID:  "room-id",
		Content: "  message body  ",
	})
	if err != nil {
		t.Fatalf("CreateChatMessage() unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("CreateChatMessage() = %#v, want %#v", got, want)
	}

	wantInput := (repository.CreateChatMessageInput{
		UserID:  "user-id",
		RoomID:  "room-id",
		Content: "message body",
	})
	if repo.createdInput != wantInput {
		t.Fatalf("repository input = %#v, want %#v", repo.createdInput, wantInput)
	}
}

func TestChatMessageUsecaseCreateChatMessageReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	repo := &stubChatMessageRepository{err: repository.ErrNotGroupMember}
	usecase := NewChatMessageUsecase(repo)

	_, err := usecase.CreateChatMessage(context.Background(), CreateChatMessageInput{
		UserID:  "user-id",
		RoomID:  "room-id",
		Content: "message",
	})
	if !errors.Is(err, repository.ErrNotGroupMember) {
		t.Fatalf("CreateChatMessage() error = %v, want %v", err, repository.ErrNotGroupMember)
	}
}
