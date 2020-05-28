package supervisor

import (
	"github.com/aporeto-inc/go-ipset/ipset"
	log "github.com/sirupsen/logrus"
)

const (
	SetBlockName = "exp-drop"
)

func BlockPerm(ip string) {
	blockList, err := ipset.New(SetBlockName, "hash:ip", &ipset.Params{})
	if err != nil {
		log.Errorf("error when blocking ip %v: %v", ip, err)
	}
	blockList.Add(ip, 0)
}
