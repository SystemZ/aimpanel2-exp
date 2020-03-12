package model

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
	err := GetOneS(&group, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return nil
	}

	return &group
}
