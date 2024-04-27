package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"keyvaluestore/cryptostore"
	"keyvaluestore/database"
	"keyvaluestore/models"
	"keyvaluestore/responses"
	"net/http"
	"time"
)

func CreateSecret(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var newSecret models.Secret
	if err := c.BindJSON(&newSecret); err != nil {
		return
	}
	newSecret.Value = string(cryptostore.EncryptText([]byte(newSecret.Value)))
	templateCollection := database.GetCollection(database.DB, "secrets")
	result, err := templateCollection.InsertOne(ctx, newSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.GenericResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, responses.GenericResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
}

func GetSecrets(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	templateCollection := database.GetCollection(database.DB, "secrets")
	result, err := templateCollection.Find(ctx, bson.D{{}})
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err := result.All(ctx, &results); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, responses.GenericResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": results}})
}

func GetSecret(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	templateCollection := database.GetCollection(database.DB, "secrets")
	var result models.Secret
	filter := bson.D{{"key", c.Param("key")}}
	err := templateCollection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		panic(err)
	}
	clear_text_string, clear_text_status := c.GetQuery("showSecret")
	if clear_text_status && clear_text_string == "True" {
		result.Value = cryptostore.DecryptText([]byte(result.Value))
	} else {
		result.Value = "*********"
	}
	c.JSON(http.StatusOK, responses.GenericResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": result}})
}

func DelSecret(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result bson.M
	filter := bson.D{{"key", c.Param("key")}}
	templateCollection := database.GetCollection(database.DB, "secrets")
	err := templateCollection.FindOneAndDelete(ctx, filter).Decode(&result)
	fmt.Println(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		panic(err)
	}
	c.JSON(http.StatusOK, responses.GenericResponse{Status: http.StatusOK, Message: "success delete", Data: map[string]interface{}{"data": result}})
}

func PutSecret(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	templateCollection := database.GetCollection(database.DB, "secrets")
	var newSecret models.Secret
	//validate the request body
	if err := c.BindJSON(&newSecret); err != nil {
		c.JSON(http.StatusBadRequest, responses.GenericResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&newSecret); validationErr != nil {
		c.JSON(http.StatusBadRequest, responses.GenericResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
		return
	}
	var result models.Secret
	filter := bson.D{{"key", c.Param("key")}}
	newSecret.Value = string(cryptostore.EncryptText([]byte(newSecret.Value)))
	err := templateCollection.FindOneAndReplace(ctx, filter, newSecret).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		panic(err)
	}

	c.JSON(http.StatusCreated, responses.GenericResponse{Status: http.StatusCreated, Message: "success replaced", Data: map[string]interface{}{"data": result}})
}
