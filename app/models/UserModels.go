package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id              string    `json:"id,omitempty"`
	Name            string    `json:"name,omitempty" validate:"required"`
	Password        string    `json:"password" validate:"required,min=8"`
	PasswordConfirm string    `json:"passwordConfirm" validate:"required"`
	Email           string    `json:"email,omitempty" validate:"required"`
	Level           string    `json:"level,omitempty" validate:"required"`
	Group           string    `json:"group,omitempty" validate:"required"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

type UserSave struct {
	Id        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty" validate:"required"`
	Password  string    `json:"password" validate:"required,min=8"`
	Email     string    `json:"email,omitempty" validate:"required"`
	Level     string    `json:"level,omitempty" validate:"required"`
	Group     string    `json:"group,omitempty" validate:"required"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type UserEdit struct {
	Id              string    `json:"id,omitempty"`
	Name            string    `json:"name,omitempty" validate:"required"`
	Password        string    `json:"password"`
	PasswordConfirm string    `json:"passwordConfirm"`
	Email           string    `json:"email,omitempty" validate:"required"`
	Level           string    `json:"level,omitempty" validate:"required"`
	Group           string    `json:"group,omitempty" validate:"required"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

type UserAuth struct {
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserAuthorization struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type DBResponse struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Id        string             `json:"id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Password  string             `json:"password" bson:"password"`
	Email     string             `json:"email" bson:"email"`
	Level     string             `json:"level,omitempty" bson:"level"`
	Group     string             `json:"group,omitempty" bson:"group"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Password  string             `json:"password" bson:"password"`
	Email     string             `json:"email" bson:"email"`
	Level     string             `json:"level,omitempty" bson:"level"`
	Group     string             `json:"group,omitempty" bson:"group"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func FilteredResponse(user *DBResponse) UserResponse {
	return UserResponse{
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type RequestResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}
