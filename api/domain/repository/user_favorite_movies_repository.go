package repository

import (
	"context"
)

type FavoriteMovieRepository interface {
	GetFavoriteMovies(ctx context.Context, userID string) ([]Movie, error)
	Create(ctx context.Context, userID string, movieID string) error
	Delete(ctx context.Context, userID string, movieID string) error
}
