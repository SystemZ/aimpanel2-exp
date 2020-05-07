package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User group represents the group for this application
// swagger:model userGroup
type GroupUser struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`
	// ID of the group
	//
	// required: true
	GroupId primitive.ObjectID `bson:"group_id" json:"group_id"`

	// ID of the user
	//
	// required: true
	UserId primitive.ObjectID `bson:"user_id" json:"user_id"`
}

func (gu *GroupUser) GetCollectionName() string {
	return groupUserCollection
}

func (gu *GroupUser) GetID() primitive.ObjectID {
	return gu.ID
}

func (gu *GroupUser) SetID(id primitive.ObjectID) {
	gu.ID = id
}

func GetGroupUserByUserId(userId primitive.ObjectID) (*GroupUser, error) {
	var gu GroupUser

	err := DB.Collection(groupUserCollection).FindOne(context.TODO(), bson.D{{Key: "user_id", Value: userId}}).Decode(&gu)
	if err != nil {
		return nil, err
	}

	return &gu, nil
}
