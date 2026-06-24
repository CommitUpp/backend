package group

import (
	"context"
	"errors"

	"github.com/CommitUpp/backend/api/domain/repository"
)

var (
	ErrUserIDRequired  = errors.New("user ID is required")
	ErrGroupIDRequired = errors.New("group ID is required")
	ErrForbidden       = errors.New("forbidden")
)

type GroupWatchedMovieUsecase interface {
	GetWatchedMovies(ctx context.Context, userID string, groupID string) ([]repository.GroupWatchedMovieRow, error)
}

type groupWatchedMovieUsecase struct {
	groupRepo             repository.GroupRepository
	groupWatchedMovieRepo repository.GroupWatchedMovieRepository
}

func NewGroupWatchedMovieUsecase(
	groupRepo repository.GroupRepository,
	groupWatchedMovieRepo repository.GroupWatchedMovieRepository,
) GroupWatchedMovieUsecase {
	return &groupWatchedMovieUsecase{
		groupRepo:             groupRepo,
		groupWatchedMovieRepo: groupWatchedMovieRepo,
	}
}

func (u *groupWatchedMovieUsecase) GetWatchedMovies(
	ctx context.Context,
	userID string,
	groupID string,
) ([]repository.GroupWatchedMovieRow, error) {
	if userID == "" {
		return nil, ErrUserIDRequired
	}

	if groupID == "" {
		return nil, ErrGroupIDRequired
	}

	isMember, err := u.groupRepo.IsGroupMember(ctx, userID, groupID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}

	watchedMovies, err := u.groupWatchedMovieRepo.GetWatchedMovies(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}

	return watchedMovies, nil
}
