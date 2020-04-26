package model

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
)

// Group represents the group for this application
// swagger:model group
type Permission struct {
	Base
	// Name of the permission
	//
	// required: true
	Name string `gorm:"column:name" json:"name"`

	// ID of the permission
	//
	// required: true
	GroupId string `gorm:"column:group_id" json:"group_id"`

	Verb uint8 `gorm:"column:verb" json:"verb"`

	Endpoint string `gorm:"column:endpoint" json:"endpoint"`
}

//TODO verify it for security vulnerability
func GetPermisionsByEndpointRegex(endpoint string) []Permission {
	var permissions []Permission
	err := Get(&permissions, map[string]interface{}{
		"selector": map[string]interface{}{
			"doc_type": "permission",
			"endpoint": map[string]interface{}{
				"$regex": endpoint,
			},
		},
	})
	if err != nil {
		return nil
	}

	return permissions
}

// FIXME return errors
func CreatePermissionsForNewHost(groupId string, hostId string) {
	perm := &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Get host",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId,
	}
	err := perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Delete host",
		Verb:     lib.GetVerbByName("DELETE"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId,
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Create game server",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "List game servers by host id",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Get host metric",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/metric",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Update host",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/update",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}
}

// FIXME return errors
func CreatePermissionsForNewGameServer(groupId string, hostId string, gameServerId string) {
	perm := &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Get game server",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId,
	}
	err := perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Delete game server",
		Verb:     lib.GetVerbByName("DELETE"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId,
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Install game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId + "/install",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Start game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId + "/start",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Restart game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId + "/restart",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Stop game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId + "/stop",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Send command to game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId + "/command",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Get logs from game server",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId + "/logs",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Console feed",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId + "/console",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "File list",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId + "/file/list",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}
}

// FIXME return errors
func CreatePermissionsForNewUser(groupId string) {
	perm := &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "List hosts",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host",
	}
	err := perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Create host",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  groupId,
		Endpoint: "/v1/host",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "List game servers by user id",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/my/server",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Change password",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  groupId,
		Endpoint: "/v1/user/change_password",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "Change email",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  groupId,
		Endpoint: "/v1/user/change_email",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "User profile",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/user/profile",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Base: Base{
			DocType: "permission",
		},
		Name:     "List games",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/game",
	}
	err = perm.Put(&perm)
	if err != nil {
		logrus.Error(err)
	}
}
