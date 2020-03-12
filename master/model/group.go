package model

import "github.com/sirupsen/logrus"

// Group represents the group for this application
// swagger:model group
type Group struct {
	Base

	// Name of the group
	//
	// required: true
	Name string `json:"name"`
}

func GetGroup(name string) *Group {
	var group Group
	if err := GetOneS(&group, map[string]string{
		"name": name,
	}); err != nil {
		logrus.Error(err)
		return nil
	}

	return &group
}
