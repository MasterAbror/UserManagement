package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/RadenAbror/UserManagement/app/config"
	"github.com/RadenAbror/UserManagement/app/helpers"
	"github.com/RadenAbror/UserManagement/app/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var levelCollection *mongo.Collection = config.GetCollection(config.DB, "user_levels")

func CreateLevel() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := c.MustGet("currentUser").(*models.DBResponse)

		user, err := helpers.FindUserById(currentUser.Id)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(
					http.StatusInternalServerError,
					models.RequestResponse{
						Status:  http.StatusInternalServerError,
						Message: "error",
						Data: map[string]interface{}{
							"data": "Identitas Anda tidak valid!",
						},
					},
				)
				return
			}
		}

		if user.Level != "63dc85d1-cfed-45cd-a404-51778f377f63" {
			c.JSON(
				http.StatusInternalServerError,
				models.RequestResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data: map[string]interface{}{
						"data": "Anda tidak diijinkan mengkases modul ini!",
					},
				},
			)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var level models.Level
		defer cancel()
		p := bluemonday.UGCPolicy()

		// validate the request body
		if err := c.BindJSON(&level); err != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Use the validator library to validate required fields
		if validationErr := validate.Struct(&level); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newLevel := models.Level{
			ID:      uuid.Must(uuid.NewRandom()).String(),
			Name:    p.Sanitize(level.Name),
			Acronym: p.Sanitize(level.Acronym),
		}

		result, err := levelCollection.InsertOne(ctx, newLevel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, models.RequestResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func DataLevel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var levels []models.Level
		defer cancel()

		results, err := levelCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleLevel models.Level
			if err = results.Decode(&singleLevel); err != nil {
				c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			levels = append(levels, singleLevel)
		}

		c.JSON(http.StatusOK,
			models.RequestResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": levels}},
		)
	}
}

func ReadLevel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		levelId := c.Param("levelId")
		var level models.Level
		defer cancel()

		objId := levelId

		err := levelCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&level)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, models.RequestResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": level}})
	}
}

func UpdateLevel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		levelId := c.Param("levelId")
		var level models.Level
		defer cancel()
		p := bluemonday.UGCPolicy()

		currentUser := c.MustGet("currentUser").(*models.DBResponse)

		user, err := helpers.FindUserById(currentUser.Id)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(
					http.StatusInternalServerError,
					models.RequestResponse{
						Status:  http.StatusInternalServerError,
						Message: "error",
						Data: map[string]interface{}{
							"data": "Identitas Anda tidak valid!",
						},
					},
				)
				return
			}
		}

		if user.Level != "63dc85d1-cfed-45cd-a404-51778f377f63" {
			c.JSON(
				http.StatusInternalServerError,
				models.RequestResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data: map[string]interface{}{
						"data": "Anda tidak diijinkan mengkases modul ini!",
					},
				},
			)
			return
		}

		//validate the request body
		if err := c.BindJSON(&level); err != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&level); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		objId := levelId
		// err := userCollection.FindOne(ctx, bson.M{"id": objId, "email": bson.M{"$ne": p.Sanitize(userEdit.Email)}}).Decode(&userEdit)

		update := bson.M{
			"name":    p.Sanitize(level.Name),
			"acronym": p.Sanitize(level.Acronym),
		}
		result, _ := levelCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

		//get updated user details
		var updatedLevel models.Level
		if result.MatchedCount == 1 {
			err := levelCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&level)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, models.RequestResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedLevel}})
	}
}

func DeleteLevel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		levelId := c.Param("levelId")
		defer cancel()

		currentUser := c.MustGet("currentUser").(*models.DBResponse)

		user, err := helpers.FindUserById(currentUser.Id)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(
					http.StatusInternalServerError,
					models.RequestResponse{
						Status:  http.StatusInternalServerError,
						Message: "error",
						Data: map[string]interface{}{
							"data": "Identitas Anda tidak valid!",
						},
					},
				)
				return
			}
		}

		if user.Level != "63dc85d1-cfed-45cd-a404-51778f377f63" {
			c.JSON(
				http.StatusInternalServerError,
				models.RequestResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data: map[string]interface{}{
						"data": "Anda tidak diijinkan mengkases modul ini!",
					},
				},
			)
			return
		}

		objId := levelId

		result, err := levelCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				models.RequestResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Level with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			models.RequestResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Level successfully deleted!"}},
		)
	}
}
