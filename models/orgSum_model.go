package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Member struct {
	UserId      string `json:"userId,omitempty" bson:"userId,omitempty"`
	PictureURL  string `json:"pictureUrl,omitempty" bson:"pictureUrl,omitempty"`
	DisplayName string `json:"displayName,omitempty" bson:"displayName,omitempty"`
	// Add other fields specific to each member object
}

type PowerUser struct {
	UserId      string `json:"userId,omitempty" bson:"userId,omitempty"`
	PictureURL  string `json:"pictureUrl,omitempty" bson:"pictureUrl,omitempty"`
	DisplayName string `json:"displayName,omitempty" bson:"displayName,omitempty"`
	// Add other fields specific to each member object
}

type Verifier struct {
	UserId      string `json:"userId,omitempty" bson:"userId,omitempty"`
	PictureURL  string `json:"pictureUrl,omitempty" bson:"pictureUrl,omitempty"`
	DisplayName string `json:"displayName,omitempty" bson:"displayName,omitempty"`
	// Add other fields specific to each member object
}

// type Image struct {}

type IdeaTemplate struct {
	Before string `json:"Before,omitempty" bson:"Before,omitempty"`
	// Image    []Image `json:"image,omitempty" bson:"image,omitempty"`
	Progress string `json:"Progress,omitempty" bson:"Progress,omitempty"`
	Learning string `json:"Learning,omitempty" bson:"Learning,omitempty"`
	Benefit  int32  `json:"benefit,omitempty" bson:"benefit,omitempty"`
}
type StatusStructure struct {
	Onhold    *int `json:"on hold,omitempty" bson:"on hold,omitempty"`
	Complete  *int `json:"complete,omitempty" bson:"complete,omitempty"`
	Ongoing   *int `json:"ongoing,omitempty" bson:"ongoing,omitempty"`
	Verifying *int `json:"verifying,omitempty" bson:"verify,omitempty"`
}

type MyOrg struct {
	AllOrgs `bson:",inline"`
	// Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	// Picture string             `json:"picture,omitempty" bson:"picture,omitempty"`
	// Member  []Member           `json:"members,omitempty" bson:"members,omitempty"`
	Share bool `json:"share,omitempty" bson:"share,omitempty"`
	// Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
}
type AllOrgs struct {
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Picture *string            `json:"picture,omitempty" bson:"picture,omitempty"`
	Member  []Member           `json:"members,omitempty" bson:"members,omitempty"`
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
}
type CreateOrg struct {
	Name      string             `json:"name,omitempty" validateOrg:"require"`
	Picture   string             `json:"picture,omitempty"`
	Member    []Member           `json:"members,omitempty" bson:"members,omitempty"`
	PowerUser []PowerUser        `json:"powerUser,omitempty" bson:"powerUser,omitempty"`
	Verifier  []Verifier         `json:"verifier,omitempty" bson:"verifier,omitempty"`
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
}
type OrgSummary struct {
	Counts       int            `json:"counts,omitempty" bson:"counts,omitempty"`
	Title        string         `json:"title,omitempty" bson:"title,omitempty"`
	Status       string         `json:"status,omitempty" bson:"status,omitempty"`
	IdeaTemplate []IdeaTemplate `json:"ideaTemplate,omitempty" bson:"ideaTemplate,omitempty"`
}
type StatusCount struct {
	Status  string `json:"status"`
	Count   int    `json:"count"`
	Benefit int32  `json:"benefit"`
}
