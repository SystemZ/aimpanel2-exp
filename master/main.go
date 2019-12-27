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
	"gitlab.com/systemz/aimpanel2/master/cmd"
)

func main() {
	cmd.Execute()
}
