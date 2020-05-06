package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User group represents the group for this application
// swagger:model userGroup
type GroupUser struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`
	// ID of the group
	//
	// required: true
	GroupId primitive.ObjectID `json:"group_id"`

	// ID of the user
	//
	// required: true
	UserId primitive.ObjectID `json:"user_id"`
}

func (g *GroupUser) GetCollectionName() string {
	return "users_group"
}

func GetGroupUserByUserId(userId primitive.ObjectID) *GroupUser {
	var gu GroupUser
	err := GetOneS(&gu, map[string]interface{}{
		"doc_type": "group_user",
		"user_id":  userId,
	})
	if err != nil {
		return nil
	}

	return &gu
}
