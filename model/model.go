package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Movie struct {
	ID       primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Title    string             `json:"title,omitempty" bson:"title,omitempty"`
	Director string             `json:"director,omitempty" bson:"director,omitempty"`
	Year     string             `json:"year,omitempty" bson:"year,omitempty"`
	Genre    string             `json:"genre,omitempty" bson:"genre,omitempty"`
	MovieID  string             `json:"movieid,omitempty" bson:"movieid,omitempty"`
}
