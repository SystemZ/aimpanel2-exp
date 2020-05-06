package model

import "go.mongodb.org/mongo-driver/bson/primitive"

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

func GetGroup(name string) *Group {
	var group Group
	err := GetOneS(&group, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return nil
	}

	return &group
}
