package models

type Level struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty" bson:"name" validate:"required"`
	Acronym string `json:"acronym,omitempty" bson:"acronym" validate:"required"`
}
