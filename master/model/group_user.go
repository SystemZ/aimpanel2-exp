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

func GetGroupUserByUserId(userId string) *GroupUser {
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
