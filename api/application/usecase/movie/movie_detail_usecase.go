package movie

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

type MovieDetailUsecase struct {
	movieDetailRepository repository.MovieDetailRepository
	groupRepository       repository.GroupRepository
}

func NewMovieDetailUsecase(
	movieDetailRepository repository.MovieDetailRepository,
	groupRepository repository.GroupRepository,
) *MovieDetailUsecase {
	return &MovieDetailUsecase{
		movieDetailRepository: movieDetailRepository,
		groupRepository:       groupRepository,
	}
}

func (u *MovieDetailUsecase) GetMovieDetail(
	ctx context.Context,
	movieID string,
	groupID string,
	userID string,
) (*repository.MovieDetail, error) {
	if userID == "" {
		return nil, ErrUserIDRequired
	}

	if groupID == "" {
		return nil, ErrGroupIDRequired
	}

	isMember, err := u.groupRepository.IsGroupMember(ctx, userID, groupID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}

	return u.movieDetailRepository.GetMovieDetail(
		ctx,
		movieID,
		groupID,
		userID,
	)
}
