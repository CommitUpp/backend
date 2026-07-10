package group

import (
	"context"

	"github.com/CommitUpp/backend/api/domain/repository"
)

type ChatRoomUsecase interface {
	GetChatRooms(ctx context.Context, userID string, groupID string, accessToken string) ([]repository.ChatRoom, error)
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
