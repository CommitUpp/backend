package movie

import (
	"context"

	"github.com/CommitUpp/backend/api/domain/repository"
)

type MoviesUsecase interface {
	GetMovies(ctx context.Context, keyword string) ([]repository.Movie, error)
}

type moviesUsecase struct {
	movieRepo repository.MovieRepository
}

func NewMoviesUsecase(repo repository.MovieRepository) MoviesUsecase {
	return &moviesUsecase{
		movieRepo: repo,
	}
}

func (u *moviesUsecase) GetMovies(ctx context.Context, keyword string) ([]repository.Movie, error) {
	return u.movieRepo.GetMovies(ctx, keyword)
}
