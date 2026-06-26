package handler

type Server struct {
	*AuthHandler
	*MoviesHandler
	*MovieDetailHandler
	*MovieStatusHandler
	*GroupHandler
	*GroupWatchedMovieHandler
}

func NewServer(
	authH *AuthHandler,
	movieH *MoviesHandler,
	movieDetailH *MovieDetailHandler,
	movieStatusH *MovieStatusHandler,
	groupH *GroupHandler,
	groupWatchedMovieH *GroupWatchedMovieHandler,
) *Server {
	return &Server{
		AuthHandler:              authH,
		MoviesHandler:            movieH,
		MovieDetailHandler:       movieDetailH,
		MovieStatusHandler:       movieStatusH,
		GroupHandler:             groupH,
		GroupWatchedMovieHandler: groupWatchedMovieH,
	}
}
