package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RadenAbror/UserManagement/app/config"
	"github.com/RadenAbror/UserManagement/app/helpers"
	"github.com/RadenAbror/UserManagement/app/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = config.GetCollection(config.DB, "users")
var validate = validator.New()

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()
		p := bluemonday.UGCPolicy()

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		if user.Password != user.PasswordConfirm {
			c.JSON(http.StatusBadRequest, models.RequestResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Konfirmasi Password tidak sama dengan Password"}})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&user)

		if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Alamat email sudah terdaftar"}})
			return
		}

		hashedPassword, _ := config.HashPassword(p.Sanitize(user.Password))
		user.Password = hashedPassword
		user.CreatedAt = time.Now()
		user.UpdatedAt = user.CreatedAt

		newUser := models.UserSave{
			Id:        primitive.NewObjectID(),
			Name:      p.Sanitize(user.Name),
			Email:     p.Sanitize(user.Email),
			Level:     user.Level,
			Group:     user.Group,
			Password:  user.Password,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, models.RequestResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userAuth models.UserAuth
		p := bluemonday.UGCPolicy()

		// validate the request body
		if err := c.BindJSON(&userAuth); err != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// use the validator library to validate required fields
		if validationErr := validate.Struct(&userAuth); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		// check email
		user, err := helpers.FindUserByEmail(p.Sanitize(userAuth.Email))
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(
					http.StatusInternalServerError,
					models.RequestResponse{
						Status:  http.StatusInternalServerError,
						Message: "error",
						Data: map[string]interface{}{
							"data": "Email tidak terdaftar!",
						},
					},
				)
				return
			}
		}

		// check password
		if err := config.VerifyPassword(user.Password, p.Sanitize(userAuth.Password)); err != nil {
			c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Password!"}})
			return
		}

		config, _ := config.LoadConfig(".")

		// Generate Tokens
		access_token, err := helpers.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		refresh_token, err := helpers.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
		c.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
		c.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

		c.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
	}
}

func GetAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		var user models.UserSave
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, models.RequestResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}})
	}
}

func EditAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		var userEdit models.UserEdit
		defer cancel()
		p := bluemonday.UGCPolicy()

		//validate the request body
		if err := c.BindJSON(&userEdit); err != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&userEdit); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.RequestResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		if userEdit.Password != userEdit.PasswordConfirm {
			c.JSON(http.StatusBadRequest, models.RequestResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Konfirmasi Password tidak sama dengan Password"}})
			return
		}

		objId, _ := primitive.ObjectIDFromHex(userId)
		// err := userCollection.FindOne(ctx, bson.M{"id": objId, "email": bson.M{"$ne": p.Sanitize(userEdit.Email)}}).Decode(&userEdit)
		count, err := userCollection.CountDocuments(
			ctx,
			bson.D{{Key: "id", Value: objId}, {Key: "email", Value: p.Sanitize(userEdit.Email)}})
		if err != nil {
			log.Fatal(err)
		}
		if count == 0 {
			checkEmail := userCollection.FindOne(ctx, bson.M{"email": p.Sanitize(userEdit.Email)}).Decode(&userEdit)
			if checkEmail != mongo.ErrNoDocuments {
				c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Email sudah digunakan!"}})
				return
			}
		}

		update := bson.M{
			"name":       p.Sanitize(userEdit.Name),
			"email":      p.Sanitize(userEdit.Email),
			"level":      userEdit.Level,
			"group":      userEdit.Group,
			"updated_at": time.Now(),
		}
		result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

		if userEdit.Password != "" {
			hashedPassword, _ := config.HashPassword(p.Sanitize(userEdit.Password))
			userEdit.Password = hashedPassword
			userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": bson.M{"password": userEdit.Password}})
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated user details
		var updatedUser models.UserSave
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, models.RequestResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}})
	}
}

func GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := c.MustGet("currentUser").(*models.DBResponse)

		c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(currentUser)}})
	}
}

func DeleteAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				models.RequestResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,

			models.RequestResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}},
		)
	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users []models.User
		defer cancel()

		results, err := userCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.User
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError, models.RequestResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			users = append(users, singleUser)
		}

		c.JSON(http.StatusOK,
			models.RequestResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": users}},
		)
	}

}

func RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := config.LoadConfig(".")

	sub, err := helpers.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := helpers.FindUserById(fmt.Sprint(sub))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	access_token, err := helpers.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
