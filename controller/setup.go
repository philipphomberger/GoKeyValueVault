package controller

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"keyvaluestore/database"
)

var validate = validator.New()

func GetCollection(name string) *mongo.Collection {
	return database.GetCollection(database.DB, name)
}
