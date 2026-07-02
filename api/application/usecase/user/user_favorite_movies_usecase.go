package user

import (
	"context"
	"github.com/CommitUpp/backend/api/domain/repository"
)

type FavoriteMovieUsecase interface {
	GetFavoriteMovies(ctx context.Context, userID string) ([]repository.Movie, error)
	CreateFavoriteMovie(ctx context.Context, userID string, movieID string) error
	DeleteFavoriteMovie(ctx context.Context, userID string, movieID string) error
}

type FavoriteMoviesUsecase struct {
	favoriteMovieRepo repository.FavoriteMovieRepository
}

func NewFavoriteMoviesUsecase(repo repository.FavoriteMovieRepository) FavoriteMovieUsecase {
	return &FavoriteMoviesUsecase{
		favoriteMovieRepo: repo,
	}
}

func (u *FavoriteMoviesUsecase) GetFavoriteMovies(
	ctx context.Context,
	userID string,
) ([]repository.Movie, error) {
	return u.favoriteMovieRepo.GetFavoriteMovies(ctx, userID)
}

func (u *FavoriteMoviesUsecase) CreateFavoriteMovie(
	ctx context.Context,
	userID string,
	movieID string,
) error {
	return u.favoriteMovieRepo.Create(ctx, userID, movieID)
}

func (u *FavoriteMoviesUsecase) DeleteFavoriteMovie(
	ctx context.Context,
	userID string,
	movieID string,
) error {
	return u.favoriteMovieRepo.Delete(ctx, userID, movieID)
}
