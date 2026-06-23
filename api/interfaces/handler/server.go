package handler

type Server struct {
	*AuthHandler
	*MoviesHandler
	*MovieStatusHandler
	*GroupHandler
	*GroupWatchedMovieHandler
}

func NewServer(authH *AuthHandler, movieH *MoviesHandler, movieStatusH *MovieStatusHandler, groupH *GroupHandler, groupWatchedMovieH *GroupWatchedMovieHandler) *Server {
	return &Server{
		AuthHandler:        authH,
		MoviesHandler:      movieH,
		MovieStatusHandler: movieStatusH,
		GroupHandler:       groupH,
		GroupWatchedMovieHandler: groupWatchedMovieH,
	}
}
