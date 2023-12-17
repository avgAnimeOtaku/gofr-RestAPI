package mongo

import (
	"context"
	"fmt"
	"log"

	"github.com/avgAnimeOtaku/gofr-RestAPI/model"
	"github.com/naamancurtis/mongo-go-struct-to-bson/mapper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gofr.dev/pkg/errors"
)

const connectionString = "mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.1.0"
const dbName = "movieDB"
const collectionName = "movies"

var collection *mongo.Collection

func init() {
	clientOption := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection success")

	collection = client.Database(dbName).Collection(collectionName)
}

func GetAllMovies() ([]primitive.M, error) {
	var movies []primitive.M

	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var movie primitive.M
		err := cursor.Decode(&movie)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	defer cursor.Close(context.Background())
	return movies, nil
}

func GetMoviesByDirector(director string) ([]primitive.M, error) {
	filter := bson.D{{Key: "director", Value: director}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error querying MongoDB:", err)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(context.Background()); err != nil {
			log.Println("Error closing cursor:", err)
		}
	}()

	var movies []primitive.M
	for cursor.Next(context.Background()) {
		var movie primitive.M
		err := cursor.Decode(&movie)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	if len(movies) == 0 {
		return nil, &errors.Response{
			Reason: "There is no movie present with the given director name",
		}
	}
	return movies, nil
}

// * Let's get movies by year
func GetMoviesByYear(year string) ([]primitive.M, error) {
	filter := bson.D{{Key: "year", Value: year}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error querying MongoDB:", err)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(context.Background()); err != nil {
			log.Println("Error closing cursor:", err)
		}
	}()

	var movies []primitive.M
	for cursor.Next(context.Background()) {
		var movie primitive.M
		err := cursor.Decode(&movie)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	if len(movies) == 0 {
		return nil, &errors.Response{
			Reason: "There are no movies present for the given year",
		}
	}
	return movies, nil
}

// * Let's get movies by genre
func GetMoviesByGenre(genre string) ([]primitive.M, error) {
	filter := bson.D{{Key: "genre", Value: genre}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error querying MongoDB:", err)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(context.Background()); err != nil {
			log.Println("Error closing cursor:", err)
		}
	}()

	var movies []primitive.M
	for cursor.Next(context.Background()) {
		var movie primitive.M
		err := cursor.Decode(&movie)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	if len(movies) == 0 {
		return nil, &errors.Response{
			Reason: "There are no movies present for the given genre",
		}
	}
	return movies, nil
}

func GetMovieByID(movieid string) (primitive.M, error) {
	filter := bson.D{{Key: "movieid", Value: movieid}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error querying MongoDB:", err)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(context.Background()); err != nil {
			log.Println("Error closing cursor:", err)
		}
	}()

	var movie primitive.M
	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&movie); err != nil {
			log.Println("Error decoding document:", err)
			return nil, err
		}
	}

	if movie == nil {
		return nil, &errors.Response{
			Reason: "There are no movies present for the given movieid",
		}
	}
	return movie, nil
}

func InsertMovie(movie model.Movie) error {
	inserted, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		return err
	}
	fmt.Println("Inserted 1 Movie in db with id: ", inserted.InsertedID)
	return nil
}

func UpdateMovie(movieid string, updateItems model.Movie) (primitive.M, error) {
	newMovie := mapper.ConvertStructToBSONMap(updateItems, nil)

	update := bson.M{"$set": newMovie}
	filter := bson.M{"movieid": movieid}

	result, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	fmt.Println("Total number of values updated are: ", result.ModifiedCount)

	movie, err := GetMovieByID(movieid)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func DeleteMovie(movieid string) (*mongo.DeleteResult, error) {

	filter := bson.M{"movieid": movieid}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	return deleteCount, nil
}
