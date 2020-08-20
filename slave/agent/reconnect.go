package agent

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type UnlimitedRetry struct {
	Verbose bool
	Name    string
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func (b UnlimitedRetry) NextBackOff() time.Duration {
	waitDuration := time.Millisecond * time.Duration(RandInt(200, 6000))
	if b.Verbose {
		logrus.Warnf("backoff %v for %v", b.Name, waitDuration.String())
	}
	return waitDuration
}
func (b *UnlimitedRetry) Reset() {}

func NewUnlimitedRetry(verbose bool, name string) *UnlimitedRetry {
	return &UnlimitedRetry{
		Verbose: verbose,
		Name:    name,
	}
}
