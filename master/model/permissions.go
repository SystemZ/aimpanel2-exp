package model

import (
	"gitlab.com/systemz/aimpanel2/lib"
	"time"
)

// Group represents the group for this application
// swagger:model group
type Permission struct {
	// ID of the group
	//
	// required: true
	ID uint `gorm:"primary_key" json:"id"`

	// Name of the permission
	//
	// required: true
	Name string `gorm:"column:name" json:"name"`

	// ID of the permission
	//
	// required: true
	GroupId int `gorm:"column:group_id" json:"group_id"`

	Verb uint8 `gorm:"column:verb" json:"verb"`

	Endpoint string `gorm:"column:endpoint" json:"endpoint"`

	// Created at timestamp
	CreatedAt time.Time `json:"created_at"`

	// Updated at timestamp
	UpdatedAt time.Time `json:"-"`

	// Deleted at timestamp
	DeletedAt *time.Time `json:"-"`
}

// FIXME return errors
func CreatePermissionsForNewHost(groupId int, hostId string) {
	DB.Save(&Permission{
		Name:     "Get host",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId,
	})

	DB.Save(&Permission{
		Name:     "Delete host",
		Verb:     lib.GetVerbByName("DELETE"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId,
	})

	DB.Save(&Permission{
		Name:     "Create game server",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server",
	})

	DB.Save(&Permission{
		Name:     "List game servers by host id",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server",
	})

	DB.Save(&Permission{
		Name:     "Get host metric",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/metric",
	})

	DB.Save(&Permission{
		Name:     "Update host",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/update",
	})
}

// FIXME return errors
func CreatePermissionsForNewGameServer(groupId int, hostId string, gameServerId string) {
	DB.Save(&Permission{
		Name:     "Get game server",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId,
	})

	DB.Save(&Permission{
		Name:     "Delete game server",
		Verb:     lib.GetVerbByName("DELETE"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId,
	})

	DB.Save(&Permission{
		Name:     "Install game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId + "/install",
	})

	DB.Save(&Permission{
		Name:     "Start game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId + "/start",
	})

	DB.Save(&Permission{
		Name:     "Restart game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId + "/restart",
	})

	DB.Save(&Permission{
		Name:     "Stop game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId + "/stop",
	})

	DB.Save(&Permission{
		Name:     "Send command to game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId + "/command",
	})

	DB.Save(&Permission{
		Name:     "Get logs from game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId + "/server/" + gameServerId + "/logs",
	})
}

// FIXME return errors
func CreatePermissionsForNewUser(groupId int) {
	DB.Save(&Permission{
		Name:     "List hosts",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host",
	})

	DB.Save(&Permission{
		Name:     "Create host",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  groupId,
		Endpoint: "/v1/host",
	})

	DB.Save(&Permission{
		Name:     "List game servers by user id",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/my/server",
	})

	DB.Save(&Permission{
		Name:     "Change password",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  groupId,
		Endpoint: "/v1/user/change_password",
	})

	DB.Save(&Permission{
		Name:     "Change email",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  groupId,
		Endpoint: "/v1/user/change_email",
	})

	DB.Save(&Permission{
		Name:     "User profile",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/user/profile",
	})

	DB.Save(&Permission{
		Name:     "List games",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/game",
	})
}
