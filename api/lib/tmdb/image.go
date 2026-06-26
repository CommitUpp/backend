package tmdb

const (
	ImageBaseURL = "https://image.tmdb.org/t/p"
	PosterSize   = "w600_and_h900_face"
	BackdropSize = "w1280"
)

func BuildPosterURL(path string) string {
	if path == "" {
		return ""
	}
	return ImageBaseURL + "/" + PosterSize + path
}

func BuildBackdropURL(path string) string {
	if path == "" {
		return ""
	}
	return ImageBaseURL + "/" + BackdropSize + path
}
