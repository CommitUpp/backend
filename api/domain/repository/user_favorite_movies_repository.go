package repository

import (
	"context"
	"errors"
)

var (
	ErrFavoriteMovieAlreadyExists = errors.New("favorite movie already exists")
	ErrFavoriteMovieLimitExceeded = errors.New("favorite movie limit exceeded")
	ErrFavoriteMovieNotFound      = errors.New("favorite movie not found")
)

type FavoriteMovieRepository interface {
	GetFavoriteMovies(ctx context.Context, userID string) ([]Movie, error)
	UpdateFavorites(ctx context.Context, userID string, addMovieIDs []string, deleteMovieIDs []string) error
	AddFavoriteMovie(ctx context.Context, userID string, movieID string) error
	DeleteFavoriteMovie(ctx context.Context, userID string, movieID string) error
}
