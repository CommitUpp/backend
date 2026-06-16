package movie

import (
	"context"
	"errors"

	"github.com/CommitUpp/backend/api/domain/repository"
)

type MovieStatusUsecase interface {
	Execute(ctx context.Context, movieID string, userID string, status string, accessToken string) error
}

type movieStatusUsecase struct {
	movieStatusRepo repository.MovieStatusRepository
}

func NewMovieStatusUsecase(repo repository.MovieStatusRepository) MovieStatusUsecase {
	return &movieStatusUsecase{
		movieStatusRepo: repo,
	}
}

func (u *movieStatusUsecase) Execute(ctx context.Context, movieID string, userID string, status string, accessToken string) error {
	// ステータスの値が正しいかをチェック
	if status != "watched" && status != "wanna_watch" {
		return errors.New("invalid status value")
	}

	err := u.movieStatusRepo.UpsertWatchStatus(ctx, movieID, userID, status, accessToken)
	if err != nil {
		return err
	}

	return nil
}
