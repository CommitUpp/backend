package handler

type Server struct {
	*AuthHandler
	*MoviesHandler
	*MovieDetailHandler
	*UserMovieStatusHandler
	*FavoriteMoviesHandler
	*GroupHandler
	*GroupWatchedMovieHandler
}

func NewServer(
	authH *AuthHandler,
	movieH *MoviesHandler,
	movieDetailH *MovieDetailHandler,
	userMovieStatusH *UserMovieStatusHandler,
	favoriteMoviesH *FavoriteMoviesHandler,
	groupH *GroupHandler,
	groupWatchedMovieH *GroupWatchedMovieHandler,
) *Server {
	return &Server{
		AuthHandler:              authH,
		MoviesHandler:            movieH,
		MovieDetailHandler:       movieDetailH,
		UserMovieStatusHandler:   userMovieStatusH,
		FavoriteMoviesHandler:    favoriteMoviesH,
		GroupHandler:             groupH,
		GroupWatchedMovieHandler: groupWatchedMovieH,
	}
}
