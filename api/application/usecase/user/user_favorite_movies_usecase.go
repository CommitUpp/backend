package user

import (
    "log"
    
	"context"
	"github.com/CommitUpp/backend/api/domain/repository"
)

type FavoriteMovieUsecase interface {
	GetFavoriteMovies(ctx context.Context, userID string) ([]repository.Movie, error)
	UpdateFavoriteMovies(ctx context.Context, userID string, movieIDs []string) error
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

func (u *FavoriteMoviesUsecase) UpdateFavoriteMovies(
	ctx context.Context,
	userID string,
	movieIDs []string,
) error {

	log.Printf("=== UpdateFavoriteMovies Start ===")
	log.Printf("UserID: %s", userID)
	log.Printf("Request MovieIDs: %v", movieIDs)

	current, err := u.favoriteMovieRepo.GetFavoriteMovies(ctx, userID)
	if err != nil {
		log.Printf("Failed to get favorite movies: %v", err)
		return err
	}

	log.Printf("Current Favorite Movies: %+v", current)

	currentMap := make(map[string]struct{})
	for _, m := range current {
		currentMap[m.MovieID] = struct{}{}
	}

	log.Printf("Current MovieID Map: %v", currentMap)

	newMap := make(map[string]struct{})
	for _, id := range movieIDs {
		newMap[id] = struct{}{}
	}

	log.Printf("New MovieID Map: %v", newMap)

	var addMovieIDs []string
	var deleteMovieIDs []string

	// 追加対象
	for id := range newMap {
		if _, ok := currentMap[id]; !ok {
			addMovieIDs = append(addMovieIDs, id)
			log.Printf("Add Target: %s", id)
		}
	}

	// 削除対象
	for id := range currentMap {
		if _, ok := newMap[id]; !ok {
			deleteMovieIDs = append(deleteMovieIDs, id)
			log.Printf("Delete Target: %s", id)
		}
	}

	log.Printf("Add MovieIDs: %v", addMovieIDs)
	log.Printf("Delete MovieIDs: %v", deleteMovieIDs)

	err = u.favoriteMovieRepo.UpdateFavorites(
		ctx,
		userID,
		addMovieIDs,
		deleteMovieIDs,
	)
	if err != nil {
		log.Printf("Failed to update favorites: %v", err)
		return err
	}

	log.Printf("Favorite movies updated successfully")
	log.Printf("=== UpdateFavoriteMovies End ===")

	return nil
}
