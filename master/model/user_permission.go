package model

import (
	"context"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/perm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type UserPermission struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	UserId primitive.ObjectID `bson:"user_id" json:"user_id" example:"1238206236281802752"`

	PermId int `bson:"perm_id" json:"perm_id"`

	HostId primitive.ObjectID `bson:"host_id" json:"host_id" example:"1238206236281802752"`

	GameServerId primitive.ObjectID `bson:"game_server_id" json:"game_server_id"`

	HostJobId primitive.ObjectID `bson:"host_job_id" json:"host_job_id" example:"1238206236281802752"`
}

func (p *UserPermission) GetCollectionName() string {
	return permissionCollection
}

func (p *UserPermission) GetID() primitive.ObjectID {
	return p.ID
}

func (p *UserPermission) SetID(id primitive.ObjectID) {
	p.ID = id
}

//TODO verify it for security vulnerability
func GetPermisionsByEndpointRegex(endpoint string) ([]UserPermission, error) {
	var permissions []UserPermission

	//cur, err := DB.Collection(permissionCollection).Find(context.TODO(),
	//	bson.D{{Key: "endpoint", Value: primitive.Regex{Pattern: endpoint, Options: ""}}})
	//if err != nil {
	//	return nil, err
	//}
	//defer cur.Close(context.TODO())
	//
	//for cur.Next(context.TODO()) {
	//	var perm Permission
	//	if err := cur.Decode(&perm); err != nil {
	//		return nil, err
	//	}
	//	permissions = append(permissions, perm)
	//}

	return permissions, nil
}

func CheckIfUserHasAccess(path string, method string, pathTemplate string, userId primitive.ObjectID) bool {
	permId, p := perm.GetByUrlAndMethod(pathTemplate, method)

	var userPerm UserPermission
	err := DB.Collection(permissionCollection).FindOne(context.TODO(), bson.D{
		{Key: "user_id", Value: userId},
		{Key: "perm_id", Value: permId},
	}).Decode(&userPerm)
	if err != nil {
		logrus.Error(err)
		return false
	}

	if p.URL == pathTemplate && p.Method == method {
		p.URL = strings.ReplaceAll(p.URL, "{hostId}", userPerm.HostId.Hex())
		p.URL = strings.ReplaceAll(p.URL, "{gsId}", userPerm.GameServerId.Hex())
		p.URL = strings.ReplaceAll(p.URL, "{jobId}", userPerm.HostJobId.Hex())

		if p.URL == path {
			return true
		}
	}

	return false
}

// FIXME return errors
func CreatePermissionsForNewHost(userId primitive.ObjectID, hostId primitive.ObjectID) error {
	perms := perm.GetKeysByGroup(perm.GroupNewHost)
	for _, p := range perms {
		userPerm := &UserPermission{
			UserId: userId,
			PermId: p,
			HostId: hostId,
		}
		err := Put(userPerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreatePermissionsForNewGameServer(userId primitive.ObjectID, hostId primitive.ObjectID, gameServerId primitive.ObjectID) error {
	perms := perm.GetKeysByGroup(perm.GroupNewGameServer)
	for _, p := range perms {
		userPerm := &UserPermission{
			UserId:       userId,
			PermId:       p,
			HostId:       hostId,
			GameServerId: gameServerId,
		}
		err := Put(userPerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreatePermissionsForNewUser(userId primitive.ObjectID) error {
	perms := perm.GetKeysByGroup(perm.GroupNewUser)
	for _, p := range perms {
		userPerm := &UserPermission{
			UserId: userId,
			PermId: p,
		}
		err := Put(userPerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreatePermissionsForNewHostJob(userId primitive.ObjectID, hostId primitive.ObjectID, jobId primitive.ObjectID) error {
	perms := perm.GetKeysByGroup(perm.GroupNewHostJob)
	for _, p := range perms {
		userPerm := &UserPermission{
			UserId:    userId,
			PermId:    p,
			HostId:    hostId,
			HostJobId: jobId,
		}
		err := Put(userPerm)
		if err != nil {
			return err
		}
	}

	return nil
}
