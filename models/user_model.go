package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId      string             `json:"userId,omitempty" bson:"userId,omitempty"`
	DisplayName string             `json:"displayName,omitempty" bson:"displayName,omitempty"`
	Picture     *interface{}       `json:"picture,omitempty" bson:"picture,omitempty"`
	PictureURL  string             `json:"pictureUrl,omitempty" bson:"pictureUrl,omitempty"`
	Profile     Profile            `json:"profile,omitempty" bson:"profile,omitempty"`
	Industry    string             `json:"industry,omitempty" bson:"industry,omitempty"`
	Province    string             `json:"province,omitempty" bson:"province,omitempty"`
}

type Profile struct {
	Language     string  `json:"language,omitempty" bson:"language,omitempty"`
	Organization *string `json:"organization,omitempty" bson:"organization,omitempty"`
}

type AccessToken struct {
	User         `bson:",inline"` // Embedded User struct
	Expire       time.Time        `json:"expire,omitempty" bson:"expire,omitempty"`
	Access_token string           `json:"access_token,omitempty" bson:"access_token,omitempty"`
}
