package repository

import (
	"context"
)

type FavoriteMovieRepository interface {
	GetFavoriteMovies(ctx context.Context, userID string) ([]Movie, error)
	UpdateFavorites(ctx context.Context, userID string, addMovieIDs []string, deleteMovieIDs []string,) error
}
