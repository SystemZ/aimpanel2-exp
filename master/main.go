// @title Aimpanel Master API
// @version 0.0.1
// @host localhost:8000
// @BasePath /v1
// @termsOfService http://swagger.io/terms/
// @securityDefinitions.apikey ApiKey
// @in header
// @name Authorization
package main

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/cmd"
	"gitlab.com/systemz/aimpanel2/master/config"
)

func main() {
	logrus.SetLevel(logrus.Level(config.LOG_LEVEL))
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    !config.DEV_MODE,
		DisableTimestamp: false,
		FullTimestamp:    true,
	})
	cmd.Execute()
}
