package migrations

import (
	"gitlab.com/systemz/aimpanel2/master/model"
)

func Migration5Up() (err error) {
	hosts, err := model.GetHosts()
	if err != nil {
		return
	}

	for _, host := range hosts {
		host.MetricMaxS = 30 * 24 * 3600
		err = model.Update(&host)
		if err != nil {
			return
		}
	}

	gameServers, err := model.GetGameServers()
	if err != nil {
		return err
	}

	for _, gs := range gameServers {
		gs.MetricMaxS = 12 * 3600
		err = model.Update(&gs)
		if err != nil {
			return
		}
	}

	return
}

func Migration5Down() {

}
