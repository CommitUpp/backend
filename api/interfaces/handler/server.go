package handler

type Server struct {
	*AuthHandler
	*MovieStatusHandler
	*GroupHandler
	*GroupWatchedMovieHandler
}

func NewServer(authH *AuthHandler, movieH *MovieStatusHandler, groupH *GroupHandler, groupWatchedMovieH *GroupWatchedMovieHandler) *Server {
	return &Server{
		AuthHandler:        authH,
		MovieStatusHandler: movieH,
		GroupHandler:       groupH,
		GroupWatchedMovieHandler: groupWatchedMovieH,
	}
}
