package lib

const (
	START   = 0
	STOP    = 1
	REBOOT  = 2
	COMMAND = 3
)

type WrapperRPC struct {
	Type int
	Body string
}
