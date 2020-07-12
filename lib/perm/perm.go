package perm

const (
	GroupNewUser       = iota
	GroupNewHost       = iota
	GroupNewGameServer = iota
	GroupNewHostJob    = iota
)

type Perm struct {
	URL    string
	Method string
	Group  uint
}

var Perms = map[int]Perm{
	//New user
	1: {URL: "/v1/host", Method: "GET", Group: GroupNewUser},
	2: {URL: "/v1/host", Method: "POST", Group: GroupNewUser},
	3: {URL: "/v1/host/my/server", Method: "GET", Group: GroupNewUser},
	4: {URL: "/v1/user/change_password", Method: "POST", Group: GroupNewUser},
	5: {URL: "/v1/user/change_email", Method: "POST", Group: GroupNewUser},
	6: {URL: "/v1/user/profile", Method: "GET", Group: GroupNewUser},
	7: {URL: "/v1/game", Method: "GET", Group: GroupNewUser},
	//New host
	8:  {URL: "/v1/host/{hostId}", Method: "GET", Group: GroupNewHost},
	9:  {URL: "/v1/host/{hostId}", Method: "DELETE", Group: GroupNewHost},
	10: {URL: "/v1/host/{hostId}/server", Method: "POST", Group: GroupNewHost},
	11: {URL: "/v1/host/{hostId}/server", Method: "GET", Group: GroupNewHost},
	12: {URL: "/v1/host/{hostId}/metric", Method: "GET", Group: GroupNewHost},
	13: {URL: "/v1/host/{hostId}/update", Method: "GET", Group: GroupNewHost},
	14: {URL: "/v1/host/{hostId}/job", Method: "POST", Group: GroupNewHost},
	15: {URL: "/v1/host/{hostId}/job", Method: "GET", Group: GroupNewHost},
	//New game server
	16: {URL: "/v1/host/{hostId}/server/{gsId}", Method: "GET", Group: GroupNewGameServer},
	17: {URL: "/v1/host/{hostId}/server/{gsId}", Method: "DELETE", Group: GroupNewGameServer},
	18: {URL: "/v1/host/{hostId}/server/{gsId}/install", Method: "PUT", Group: GroupNewGameServer},
	19: {URL: "/v1/host/{hostId}/server/{gsId}/start", Method: "PUT", Group: GroupNewGameServer},
	20: {URL: "/v1/host/{hostId}/server/{gsId}/restart", Method: "PUT", Group: GroupNewGameServer},
	21: {URL: "/v1/host/{hostId}/server/{gsId}/stop", Method: "PUT", Group: GroupNewGameServer},
	22: {URL: "/v1/host/{hostId}/server/{gsId}/command", Method: "PUT", Group: GroupNewGameServer},
	23: {URL: "/v1/host/{hostId}/server/{gsId}/logs", Method: "GET", Group: GroupNewGameServer},
	24: {URL: "/v1/host/{hostId}/server/{gsId}/console", Method: "GET", Group: GroupNewGameServer},
	25: {URL: "/v1/host/{hostId}/server/{gsId}/file/list", Method: "GET", Group: GroupNewGameServer},
	//Host job
	26: {URL: "/v1/host/{hostId}/job/{jobId}", Method: "DELETE", Group: GroupNewHostJob},
}

func GetByUrlAndMethod(url string, method string) (int, Perm) {
	for k, v := range Perms {
		if v.URL == url && v.Method == method {
			return k, v
		}
	}

	return 0, Perm{}
}

func GetKeysByGroup(group uint) []int {
	var keys []int

	for k, v := range Perms {
		if v.Group == group {
			keys = append(keys, k)
		}
	}

	return keys
}
