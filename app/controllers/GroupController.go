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
	"go.mongodb.org/mongo-driver/mongo"
)

var groupCollection *mongo.Collection = config.GetCollection(config.DB, "user_groups")

func CreateGroup() gin.HandlerFunc {
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
		var group models.Group
		defer cancel()
		p := bluemonday.UGCPolicy()

		// Validate the request body
		if err := c.BindJSON(&group); err != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		if validationErr := validate.Struct(&group); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data: map[string]interface{}{
					"data": validationErr.Error(),
				},
			})
		}

		newGroup := models.Group{
			ID:      uuid.Must(uuid.NewRandom()).String(),
			Name:    p.Sanitize(group.Name),
			Acronym: p.Sanitize(group.Acronym),
		}
		result, err := groupCollection.InsertOne(ctx, newGroup)

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.RequestResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
		}

		c.JSON(http.StatusCreated, models.RequestResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data: map[string]interface{}{
				"data": result,
			},
		})
	}
}

func DataGroup() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func ReadGroup() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateGroup() gin.HandlerFunc {
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
		var group models.Group
		defer cancel()
		p := bluemonday.UGCPolicy()

		// Validate the request body
		if err := c.BindJSON(&group); err != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		if validationErr := validate.Struct(&group); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data: map[string]interface{}{
					"data": validationErr.Error(),
				},
			})
		}

		newGroup := models.Group{
			ID:      uuid.Must(uuid.NewRandom()).String(),
			Name:    p.Sanitize(group.Name),
			Acronym: p.Sanitize(group.Acronym),
		}
		result, err := groupCollection.InsertOne(ctx, newGroup)

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.RequestResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
		}

		c.JSON(http.StatusCreated, models.RequestResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data: map[string]interface{}{
				"data": result,
			},
		})
	}
}

func DeleteGroup() gin.HandlerFunc {
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
	}
}
