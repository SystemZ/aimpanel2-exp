package supervisor

import "gitlab.com/systemz/aimpanel2/slave/model"

func Start() {
	//Init redis
	model.InitRedis()

	redisStarted := make(chan bool, 1)
	go listenerRedis(redisStarted)

	InstallPackages()
	ReportInit()
	WatchSshd()
}
