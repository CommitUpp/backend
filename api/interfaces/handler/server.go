package handler

type Server struct {
	*AuthHandler
	*MovieStatusHandler
	*GroupHandler
}

func NewServer(authH *AuthHandler, movieH *MovieStatusHandler, groupH *GroupHandler) *Server {
	return &Server{
		AuthHandler:        authH,
		MovieStatusHandler: movieH,
		GroupHandler:       groupH,
	}
}
