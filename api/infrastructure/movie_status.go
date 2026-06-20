package infrastructure

import (
	"context"
	"strings"
	"time"

	"github.com/CommitUpp/backend/api/domain/repository"
	"github.com/supabase-community/postgrest-go"
)

type MovieStatusRepository struct {
	baseURL string
	apiKey  string
}

func NewMovieStatusRepository(baseURL string, apiKey string) repository.MovieStatusRepository {
	return &MovieStatusRepository{
		baseURL: strings.TrimRight(baseURL, "/") + "/rest/v1",
		apiKey:  apiKey,
	}
}

func (r *MovieStatusRepository) WatchStatus(ctx context.Context, movieID string, userID string, status string, accessToken string) error {
	// データを作成
	row := map[string]interface{}{
		"user_id":  userID,
		"movie_id": movieID,
		"status":   status,
	}

	client := postgrest.NewClient(
		r.baseURL,
		"public",
		map[string]string{
			"apikey":        r.apiKey,
			"Authorization": "Bearer " + accessToken,
		},
	)
	if client.ClientError != nil {
		return client.ClientError
	}

	// watch_statusesに登録
	_, _, err := client.
		From("watch_statuses").
		Upsert(row, "user_id,movie_id", "minimal", "merge").
		ExecuteWithContext(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (r *MovieStatusRepository) GetWatchStatuses(ctx context.Context, userID string, status *string, accessToken string) ([]repository.MovieStatus, error) {
	client := postgrest.NewClient(
		r.baseURL,
		"public",
		map[string]string{
			"apikey":        r.apiKey,
			"Authorization": "Bearer " + accessToken,
		},
	)
	if client.ClientError != nil {
		return nil, client.ClientError
	}

	query := client.
		From("watch_statuses").
		Select(
			"movie_id,updated_at,movies!watch_statuses_movie_id_fkey(tmdb_id,title,poster_url,trailer_url,overview,release_date)",
			"",
			false,
		).
		Eq("user_id", userID).
		Order("updated_at", &postgrest.OrderOpts{Ascending: false, NullsFirst: false})

	if status != nil {
		query = query.Eq("status", *status)
	}

	var rows []watchStatusRow
	if _, err := query.ExecuteToWithContext(ctx, &rows); err != nil {
		return nil, err
	}

	movies := make([]repository.MovieStatus, 0, len(rows))
	for _, row := range rows {
		movies = append(movies, repository.MovieStatus{
			MovieID:     row.MovieID,
			TMDBID:      row.Movie.TMDBID,
			Title:       row.Movie.Title,
			PosterURL:   row.Movie.PosterURL,
			TrailerURL:  row.Movie.TrailerURL,
			Overview:    row.Movie.Overview,
			ReleaseDate: row.Movie.ReleaseDate,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	return movies, nil
}

type watchStatusRow struct {
	MovieID   string       `json:"movie_id"`
	UpdatedAt time.Time    `json:"updated_at"`
	Movie     movieRowData `json:"movies"`
}

type movieRowData struct {
	TMDBID      string `json:"tmdb_id"`
	Title       string `json:"title"`
	PosterURL   string `json:"poster_url"`
	TrailerURL  string `json:"trailer_url"`
	Overview    string `json:"overview"`
	ReleaseDate string `json:"release_date"`
}
