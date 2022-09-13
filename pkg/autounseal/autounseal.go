package autounseal

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	ServiceDomain         string
	ServiceScheme         string
	ServicePort           int
	UnlockShares          int
	UnlockThreshold       int
	VaultRootTokenSecret  string
	VaultUnlockKeysSecret string
}

func ManageVaultAutounseal(cfg Config) {
	ips, err := net.LookupIP(cfg.ServiceDomain)
	if err != nil {
		panic(err)
	}
	log.Infof("Found ips: %v", ips)

	// check their statuses
	// if all uninitialize
	//   delete secrets
	//   initialize
	//   store keys and root token in secrets
	// else if all initialized
	//   check if any locked, then unlock it
	time.Sleep(30 * time.Second)
}
