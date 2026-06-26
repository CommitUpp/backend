package movie

import (
	"context"

	"github.com/CommitUpp/backend/api/domain/repository"
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
) (*repository.MovieDetail, error) {
	return u.movieDetailRepository.GetMovieDetail(ctx, movieID)
}
