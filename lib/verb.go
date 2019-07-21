package lib

var verbs = map[string]uint8{
	"GET":    1,
	"POST":   2,
	"PUT":    3,
	"DELETE": 4,
}

func GetVerbByName(name string) uint8 {
	return verbs[name]
}
