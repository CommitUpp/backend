package handler

type Server struct {
	*AuthHandler
	*MovieStatusHandler
}

func NewServer(authH *AuthHandler, movieH *MovieStatusHandler) *Server {
	return &Server{
		AuthHandler:        authH,
		MovieStatusHandler: movieH,
	}
}
