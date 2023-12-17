package main

import (
	"github.com/avgAnimeOtaku/gofr-RestAPI/controller"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.GET("/movies", controller.GetMovies)
	app.GET("/movies/movieid/{movieidid}", controller.GetMovieID)
	app.GET("/movies/director/{director}", controller.GetMovieDirector)
	app.GET("/movies/genre/{genre}", controller.GetMovieGenre)
	app.GET("/movies/year/{year}", controller.GetMovieYear)
	app.POST("/movies", controller.CreateMovie)
	app.PUT("/movies/{movieid}", controller.UpdateMovie)
	app.DELETE("/movies/{movieid}", controller.DeleteMovie)

	app.Start()
}
