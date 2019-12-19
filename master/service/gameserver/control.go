package gameserver

import (
	"errors"
	"github.com/alexandrevicenzi/go-sse"
	"gitlab.com/systemz/aimpanel2/lib/rabbit"
	"gitlab.com/systemz/aimpanel2/master/events"
	"gitlab.com/systemz/aimpanel2/master/model"
	"strconv"
	"time"
)

func Start(gsId string) error {
	gameServer := model.GetGameServer(model.DB, gsId)
	if gameServer == nil {
		return errors.New("error when getting game server from db")
	}

	hostToken := model.GetHostToken(model.DB, gameServer.HostId.String())
	if hostToken == "" {
		return errors.New("error when getting host token from db")
	}

	channel, ok := events.SSE.GetChannel("/v1/events/" + hostToken)
	if !ok {
		return errors.New("host is not turned on")
	}
	channel.SendMessage(sse.NewMessage("", gameServer.ID.String(), strconv.Itoa(rabbit.WRAPPER_START)))

	model.Redis.Set("gs_start_id_"+gameServer.ID.String(), 1, 1*time.Hour)

	return nil
}
