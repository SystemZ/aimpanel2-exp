package supervisor

import (
	"github.com/arduino/go-apt-client"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

var (
	packagesRefreshed = false
)

func InstallPackages() {
	// easy blacklisting massive number of IPs
	InstallPackage("ipset", "ipset")
}

// FIXME return error
func InstallPackage(nameInPath string, packageName string) {
	log.Infof("Checking if %v is installed...", packageName)
	_, err := exec.LookPath(nameInPath)
	if err != nil {
		log.Infof("%v not found, installing...", packageName)
		if !packagesRefreshed {
			_, err = apt.CheckForUpdates()
			if err != nil {
				log.Warnf("error when trying to update package list: %v", err)
			}
			packagesRefreshed = true
		}
		apt.Install(&apt.Package{Name: packageName})
	}
	_, err = exec.LookPath(nameInPath)
	if err != nil {
		log.Errorf("%v package installation failed")
		return
	}
	log.Infof("%v is installed", packageName)
}
