package model

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Group represents the group for this application
// swagger:model group
type Permission struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	// Name of the permission
	//
	// required: true
	Name string `gorm:"column:name" json:"name"`

	// ID of the permission
	//
	// required: true
	GroupId primitive.ObjectID `gorm:"column:group_id" json:"group_id"`

	Verb uint8 `gorm:"column:verb" json:"verb"`

	Endpoint string `gorm:"column:endpoint" json:"endpoint"`
}

func (p *Permission) GetCollectionName() string {
	return "permissions"
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
func CreatePermissionsForNewHost(groupId primitive.ObjectID, hostId primitive.ObjectID) {
	perm := &Permission{
		Name:     "Get host",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String(),
	}
	err := Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Delete host",
		Verb:     lib.GetVerbByName("DELETE"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String(),
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Create game server",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/server",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "List game servers by host id",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/server",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Get host metric",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/metric",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Update host",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/update",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Create host job",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/job",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Get host jobs",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/job",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}
}

// FIXME return errors
func CreatePermissionsForNewGameServer(groupId primitive.ObjectID, hostId primitive.ObjectID, gameServerId primitive.ObjectID) {
	perm := &Permission{
		Name:     "Get game server",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/server/" + gameServerId.String(),
	}
	err := Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Delete game server",
		Verb:     lib.GetVerbByName("DELETE"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/server/" + gameServerId.String(),
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Install game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/server/" + gameServerId.String() + "/install",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Start game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/server/" + gameServerId.String() + "/start",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Restart game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/server/" + gameServerId.String() + "/restart",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Stop game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/server/" + gameServerId.String() + "/stop",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Send command to game server",
		Verb:     lib.GetVerbByName("PUT"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/server/" + gameServerId.String() + "/command",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Get logs from game server",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/server/" + gameServerId.String() + "/logs",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Console feed",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/server/" + gameServerId.String() + "/console",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "File list",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/server/" + gameServerId.String() + "/file/list",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}
}

// FIXME return errors
func CreatePermissionsForNewUser(groupId primitive.ObjectID) {
	perm := &Permission{
		Name:     "List hosts",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host",
	}
	err := Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Create host",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  groupId,
		Endpoint: "/v1/host",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "List game servers by user id",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/host/my/server",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Change password",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  groupId,
		Endpoint: "/v1/user/change_password",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "Change email",
		Verb:     lib.GetVerbByName("POST"),
		GroupId:  groupId,
		Endpoint: "/v1/user/change_email",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "User profile",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/user/profile",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}

	perm = &Permission{
		Name:     "List games",
		Verb:     lib.GetVerbByName("GET"),
		GroupId:  groupId,
		Endpoint: "/v1/game",
	}
	err = Put(perm)
	if err != nil {
		logrus.Error(err)
	}
}

func CreatePermissionsForNewHostJob(groupId primitive.ObjectID, hostId primitive.ObjectID, jobId primitive.ObjectID) {
	perm := &Permission{
		Name:     "Remove host job",
		Verb:     lib.GetVerbByName("DELETE"),
		GroupId:  groupId,
		Endpoint: "/v1/host/" + hostId.String() + "/job/" + jobId.String(),
	}
	err := Put(perm)
	if err != nil {
		logrus.Error(err)
	}
}
