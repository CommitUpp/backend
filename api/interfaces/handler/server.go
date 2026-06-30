package handler

type Server struct {
	*AuthHandler
	*MoviesHandler
	*MovieDetailHandler
	*UserMovieStatusHandler
	*GroupHandler
	*GroupWatchedMovieHandler
}

func NewServer(
	authH *AuthHandler,
	movieH *MoviesHandler,
	movieDetailH *MovieDetailHandler,
	userMovieStatusH *UserMovieStatusHandler,
	groupH *GroupHandler,
	groupWatchedMovieH *GroupWatchedMovieHandler,
) *Server {
	return &Server{
		AuthHandler:              authH,
		MoviesHandler:            movieH,
		MovieDetailHandler:       movieDetailH,
		UserMovieStatusHandler:   userMovieStatusH,
		GroupHandler:             groupH,
		GroupWatchedMovieHandler: groupWatchedMovieH,
	}
}
