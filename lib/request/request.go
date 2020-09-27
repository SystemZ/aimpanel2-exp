package request

import (
	"fmt"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"strings"
)

type HostCreate struct {
	// User assigned name
	Name *string `json:"name" example:"My Great Linux server"`

	// User assigned ip
	Ip *string `json:"ip" example:"192.51.100.128"`

	HwId *string `json:"hw_id"`
}

type HostCreateJob struct {
	// User assigned name
	Name string `json:"name" example:"My Great job"`

	CronExpression string `json:"cron_expression"`

	TaskMessage task.Message `json:"task_message"`
}

type GameServerCreate struct {
	//User assigned name
	Name *string `json:"name" example:"Ultra MC Server"`

	//Selected game id
	GameId *uint `json:"game_id" example:"1"`

	//Selected game version
	GameVersion *string `json:"game_version" example:"1.14.2"`

	//Custom cmd for starting game server
	CustomCmdStart *string `json:"custom_cmd_start" example:"java -jar mc.jar"`

	//Ports used by game server
	Ports *[]GameServerPort `json:"ports"`
}

type GameServerPort struct {
	Protocol      string `json:"protocol,omitempty"`
	Host          string `json:"host,omitempty"`
	PortHost      int    `json:"port_host,omitempty"`
	PortContainer int    `json:"port_container,omitempty"`
}

func (gs GameServerCreate) SerializePorts() (res string) {
	var resPorts []string
	for _, port := range *gs.Ports {
		resPorts = append(resPorts, fmt.Sprintf("%v:%v->%v/%v", port.Host, port.PortHost, port.PortContainer, port.Protocol))
	}
	res = strings.Join(resPorts, " ")
	return
}

type GameServerStop struct {
	Type uint `json:"type" example:"1"`
}

type GameServerSendCommand struct {
	Command string `json:"command" example:"say Hello!"`
}

//swagger:parameters User changePassword
type UserChangePassword struct {
	Password          string `json:"password"`
	NewPassword       string `json:"new_password"`
	NewPasswordRepeat string `json:"new_password_repeat"`
}

//swagger:parameters User changeEmail
type UserChangeEmail struct {
	Email          string `json:"email"`
	NewEmail       string `json:"new_email"`
	NewEmailRepeat string `json:"new_email_repeat"`
}

type AuthRegister struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
	Email          string `json:"email"`
	EmailRepeat    string `json:"email_repeat"`
}

type AuthLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
