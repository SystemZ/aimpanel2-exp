package model

// User group represents the group for this application
// swagger:model userGroup
type GroupUser struct {
	Base
	// ID of the group
	//
	// required: true
	GroupId string `json:"group_id"`

	// ID of the user
	//
	// required: true
	UserId string `json:"user_id"`
}
