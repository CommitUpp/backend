package handler

type Server struct {
	*AuthHandler
	*MoviesHandler
	*MovieDetailHandler
	*UserMovieStatusHandler
	*MyFavoriteMoviesHandler
	*UserFavoriteMoviesHandler
	*GroupHandler
	*GroupWatchedMovieHandler
}

func NewServer(
	authH *AuthHandler,
	movieH *MoviesHandler,
	movieDetailH *MovieDetailHandler,
	userMovieStatusH *UserMovieStatusHandler,
	myFavoriteMoviesH *MyFavoriteMoviesHandler,
	userFavoriteMoviesH *UserFavoriteMoviesHandler,
	groupH *GroupHandler,
	groupWatchedMovieH *GroupWatchedMovieHandler,
) *Server {
	return &Server{
		AuthHandler:              authH,
		MoviesHandler:            movieH,
		MovieDetailHandler:       movieDetailH,
		UserMovieStatusHandler:   userMovieStatusH,
		MyFavoriteMoviesHandler:  myFavoriteMoviesH,
		UserFavoriteMoviesHandler: userFavoriteMoviesH,
		GroupHandler:             groupH,
		GroupWatchedMovieHandler: groupWatchedMovieH,
	}
}
