package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Person struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	FirstName   string             `json:"first_name"`
	LastName    string             `json:"last_name"`
	Address     string             `json:"address"`
	PhoneNumber string             `json:"phone_number"`
}
