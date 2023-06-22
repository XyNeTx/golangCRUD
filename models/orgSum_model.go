package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Member struct {
	UserId      string `json:"userId,omitempty" bson:"userId,omitempty"`
	PictureURL  string `json:"pictureUrl,omitempty" bson:"pictureUrl,omitempty"`
	DisplayName string `json:"displayName,omitempty" bson:"displayName,omitempty"`
	// Add other fields specific to each member object
}

type MyOrg struct {
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Picture string             `json:"picture,omitempty" bson:"picture,omitempty"`
	Member  []Member           `json:"members,omitempty" bson:"members,omitempty"`
	Share   bool               `json:"share,omitempty" bson:"share,omitempty"`
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
}
type AllOrgs struct {
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Picture string             `json:"picture,omitempty" bson:"picture,omitempty"`
	Member  []Member           `json:"members,omitempty" bson:"members,omitempty"`
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
}
type CreateOrg struct {
	Name    string             `json:"name,omitempty" validate:"require"`
	Picture string             `json:"picture,omitempty"`
	Id      primitive.ObjectID `json:"id,omitempty"`
}
type OrgSummary struct {
	Counts     primitive.A `json:"counts,omitempty" bson:"counts,omitempty"`
	Status     string      `json:"status,omitempty" bson:"status,omitempty"`
	Value      string      `json:"value,omitempty" bson:"value,omitempty"`
	TotalCount int         `json:"totalCount,omitempty" bson:"totalCount,omitempty"`
	Ideas      primitive.A `json:"ideas,omitempty" bson:"ideas,omitempty"`
	Benefit    primitive.A `json:"benefit,omitempty" bson:"benefit,omitempty"`
}
