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
}

func NewMovieDetailUsecase(
	movieDetailRepository repository.MovieDetailRepository,
) *MovieDetailUsecase {
	return &MovieDetailUsecase{
		movieDetailRepository: movieDetailRepository,
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

	isMember, err := u.movieDetailRepository.IsGroupMember(ctx, groupID, userID)
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
