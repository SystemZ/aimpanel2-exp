package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Group represents the group for this application
// swagger:model group
type Group struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	// Name of the group
	//
	// required: true
	Name string `json:"name"`
}

func (g *Group) GetCollectionName() string {
	return "groups"
}

func (g *Group) GetID() primitive.ObjectID {
	return g.ID
}

func GetGroupByName(name string) (*Group, error) {
	var group Group

	err := DB.Collection(groupCollection).FindOne(context.TODO(), bson.D{{"name", name}}).
		Decode(&group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}
