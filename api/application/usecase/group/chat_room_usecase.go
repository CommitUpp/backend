package group

import (
	"context"
	"errors"

	"github.com/CommitUpp/backend/api/domain/repository"
)

var (
	ErrMovieIDRequired       = errors.New("movie ID is required")
	ErrMovieNotFound         = repository.ErrMovieNotFound
	ErrMovieNotWatched       = repository.ErrMovieNotWatched
	ErrChatRoomAlreadyExists = repository.ErrChatRoomAlreadyExists
)

type ChatRoomUsecase interface {
	GetChatRooms(ctx context.Context, userID string, groupID string, accessToken string) ([]repository.ChatRoom, error)
	CreateChatRoom(ctx context.Context, input CreateChatRoomInput) (repository.ChatRoom, error)
}

type CreateChatRoomInput struct {
	UserID      string
	GroupID     string
	MovieID     string
	AccessToken string
}

type chatRoomUsecase struct {
	groupRepo    repository.GroupRepository
	chatRoomRepo repository.ChatRoomRepository
}

func NewChatRoomUsecase(
	groupRepo repository.GroupRepository,
	chatRoomRepo repository.ChatRoomRepository,
) ChatRoomUsecase {
	return &chatRoomUsecase{
		groupRepo:    groupRepo,
		chatRoomRepo: chatRoomRepo,
	}
}

func (u *chatRoomUsecase) GetChatRooms(
	ctx context.Context,
	userID string,
	groupID string,
	accessToken string,
) ([]repository.ChatRoom, error) {
	if userID == "" {
		return nil, ErrUserIDRequired
	}

	if groupID == "" {
		return nil, ErrGroupIDRequired
	}

	if accessToken == "" {
		return nil, ErrAccessTokenRequired
	}

	isMember, err := u.groupRepo.IsGroupMember(ctx, userID, groupID, accessToken)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}

	return u.chatRoomRepo.GetChatRooms(ctx, groupID)
}

func (u *chatRoomUsecase) CreateChatRoom(
	ctx context.Context,
	input CreateChatRoomInput,
) (repository.ChatRoom, error) {
	if input.UserID == "" {
		return repository.ChatRoom{}, ErrUserIDRequired
	}
	if input.GroupID == "" {
		return repository.ChatRoom{}, ErrGroupIDRequired
	}
	if input.MovieID == "" {
		return repository.ChatRoom{}, ErrMovieIDRequired
	}
	if input.AccessToken == "" {
		return repository.ChatRoom{}, ErrAccessTokenRequired
	}

	isMember, err := u.groupRepo.IsGroupMember(
		ctx,
		input.UserID,
		input.GroupID,
		input.AccessToken,
	)
	if err != nil {
		return repository.ChatRoom{}, err
	}
	if !isMember {
		return repository.ChatRoom{}, ErrForbidden
	}

	return u.chatRoomRepo.CreateChatRoom(ctx, repository.CreateChatRoomInput{
		GroupID:   input.GroupID,
		MovieID:   input.MovieID,
		CreatedBy: input.UserID,
	})
}
