package host

import (
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
	"os"
	"time"
)

func Create(data *request.HostCreate, userId string) (*model.Host, int) {
	host := &model.Host{
		Base: model.Base{
			DocType: "host",
		},
		Name:            data.Name,
		Ip:              data.Ip,
		UserId:          userId,
		MetricFrequency: 30,
		Token:           lib.RandomString(32),
	}
	err := host.Put(&host)
	if err != nil {
		return nil, ecode.DbSave
	}

	group := model.GetGroup("USER-" + userId)
	if group == nil {
		return nil, ecode.GroupNotFound
	}

	// FIXME handle errors
	model.CreatePermissionsForNewHost(group.ID, host.ID)

	return host, ecode.NoError
}

//Removes host and linked game servers
func Remove(hostId string) int {
	host := model.GetHost(hostId)
	gameServers := model.GetGameServersByHostId(hostId)
	for _, gameServer := range *gameServers {
		err := gameserver.Remove(gameServer.ID)
		if err != nil {
			return ecode.GsRemove
		}
	}

	permissions := model.GetPermisionsByEndpointRegex("/v1/host/" + host.ID)
	for _, perm := range permissions {
		err := model.Delete(perm.ID, perm.Rev)
		if err != nil {
			return ecode.DbError
		}
	}

	err := model.Delete(host.ID, host.Rev)
	if err != nil {
		return ecode.DbError
	}

	return ecode.NoError
}

func Auth(t string) (string, int) {
	host := model.GetHostByToken(t)

	if host == nil {
		return "", ecode.HostNotFound
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 48).Unix(),
		"uid": host.ID,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", ecode.Unknown
	}

	return tokenString, ecode.NoError
}
