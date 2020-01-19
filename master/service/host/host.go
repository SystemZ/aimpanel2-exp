package host

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/lib/request"
	"gitlab.com/systemz/aimpanel2/master/model"
	"gitlab.com/systemz/aimpanel2/master/service/gameserver"
	"os"
	"time"
)

func Create(data *request.HostCreate, userId uuid.UUID) (*model.Host, int) {
	host := &model.Host{
		Name:            data.Name,
		Ip:              data.Ip,
		UserId:          userId,
		MetricFrequency: 30,
	}
	model.DB.Save(&host)

	group := model.GetGroup(model.DB, "USER-"+userId.String())
	if group == nil {
		return nil, ecode.GroupNotFound
	}

	// FIXME handle errors
	model.CreatePermissionsForNewHost(group.ID, host.ID.String())

	return host, ecode.NoError
}

//Removes host and linked game servers
func Remove(hostId string) int {
	host := model.GetHost(model.DB, hostId)
	gameServers := model.GetGameServersByHostId(model.DB, hostId)
	for _, gameServer := range *gameServers {
		err := gameserver.Remove(gameServer.ID.String())
		if err != nil {
			return ecode.GsRemove
		}
	}

	model.DB.Where("endpoint LIKE ?", "/v1/host/"+host.ID.String()+"%").Delete(&model.Permission{})
	model.DB.Delete(&host)

	return ecode.NoError
}

func Auth(t string) (string, int) {
	host := model.GetHostByToken(model.DB, t)

	if host == nil {
		return "", ecode.HostNotFound
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 48).Unix(),
		"uid": host.ID.String(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", ecode.Unknown
	}

	return tokenString, ecode.NoError
}
