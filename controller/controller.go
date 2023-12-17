package controller

import (
	"encoding/json"
	"fmt"

	"github.com/avgAnimeOtaku/gofr-RestAPI/model"
	"github.com/avgAnimeOtaku/gofr-RestAPI/mongo"

	"gofr.dev/pkg/errors"
	"gofr.dev/pkg/gofr"
)

func GetMovies(gct *gofr.Context) (interface{}, error) {
	movies, err := mongo.GetAllMovies()
	if err != nil {
		return nil, err
	}

	var theMovies []model.Movie
	movieByte, _ := json.Marshal(movies)
	json.Unmarshal(movieByte, &theMovies)
	return theMovies, nil
}

func GetMovieID(gct *gofr.Context) (interface{}, error) {
	movieid := gct.PathParam("movieid")
	movies, err := mongo.GetMovieByID(movieid)

	if err != nil {
		return nil, err
	}

	var theMovies model.Movie
	movieByte, _ := json.Marshal(movies)
	json.Unmarshal(movieByte, &theMovies)
	return theMovies, nil
}

func GetMovieDirector(gct *gofr.Context) (interface{}, error) {
	director := gct.PathParam("director")
	fmt.Println("the director name is", director)
	movies, err := mongo.GetMoviesByDirector(director)

	if err != nil {
		return nil, err
	}

	var theMovies []model.Movie
	movieByte, _ := json.Marshal(movies)
	json.Unmarshal(movieByte, &theMovies)
	return theMovies, nil
}

func GetMovieYear(gct *gofr.Context) (interface{}, error) {
	year := gct.PathParam("year")
	fmt.Println("the year is", year)
	movies, err := mongo.GetMoviesByYear(year)

	if err != nil {
		return nil, err
	}

	var theMovie []model.Movie
	movieByte, _ := json.Marshal(movies)
	json.Unmarshal(movieByte, &theMovie)
	return theMovie, nil
}

func GetMovieGenre(gct *gofr.Context) (interface{}, error) {
	genre := gct.PathParam("genre")
	fmt.Println("the genre is", genre)
	movies, err := mongo.GetMoviesByGenre(genre)

	if err != nil {
		return nil, err
	}

	var theMovies []model.Movie
	movieByte, _ := json.Marshal(movies)
	json.Unmarshal(movieByte, &theMovies)
	return theMovies, nil
}

func CreateMovie(gct *gofr.Context) (interface{}, error) {
	var movie model.Movie
	gct.Bind(&movie)

	if err := isJsonValid(movie); err != nil {
		return nil, err
	}

	newMovieID := movie.ID.Hex()
	_, err := mongo.GetMovieByID(newMovieID)

	if err == nil {
		return nil, &errors.Response{
			Reason: "There is already a movie present with the given ID",
		}
	}

	err = mongo.InsertMovie(movie)
	if err != nil {
		return nil, err
	}

	var theMovie model.Movie
	movieByte, _ := json.Marshal(movie)
	json.Unmarshal(movieByte, &theMovie)
	return theMovie, nil
}

func UpdateMovie(gct *gofr.Context) (interface{}, error) {
	var movie model.Movie
	gct.Bind(&movie)

	id := gct.PathParam("movieid")
	if movie.MovieID != "" {
		return nil, &errors.Response{
			Reason: "ID could not be updated once set",
		}
	}

	if movie.Director == "" && movie.Title == "" && movie.Genre == "" && movie.Year == "" {
		return nil, &errors.Response{
			Reason: "Check the format of data, or the name of the fields, it should be director, title, and genre",
		}
	}

	uMovie, err := mongo.UpdateMovie(id, movie)
	if err != nil {
		return nil, err
	}

	var theMovie model.Movie
	movieByte, _ := json.Marshal(uMovie)
	json.Unmarshal(movieByte, &theMovie)
	return theMovie, nil
}

func DeleteMovie(gct *gofr.Context) (interface{}, error) {
	movieid := gct.PathParam("movieid")
	_, err := mongo.GetMovieByID(movieid)

	if err != nil {
		return nil, &errors.Response{
			Reason: "There is no movie present with the given ID",
		}
	}

	deletedMovie, err := mongo.DeleteMovie(movieid)
	if err != nil {
		return nil, err
	}
	return deletedMovie, nil
}

func isJsonValid(movie model.Movie) error {
	byteMovie, _ := json.Marshal(movie)
	valid := json.Valid(byteMovie)

	if !valid {
		return &errors.Response{
			Reason: "Json is not in valid format",
		}
	}

	if movie.Title == "" {
		return &errors.Response{
			Reason: "The title of the movie is missing",
		}
	}

	if movie.Director == "" {
		return &errors.Response{
			Reason: "The director of the movie is missing",
		}
	}

	if movie.Genre == "" {
		return &errors.Response{
			Reason: "The genre of the movie is missing",
		}
	}

	if movie.Year == "" {
		return &errors.Response{
			Reason: "The year of the movie is missing",
		}
	}

	if movie.MovieID == "" {
		return &errors.Response{
			Reason: "The movieid of the movie is missing",
		}
	}

	return nil
}
