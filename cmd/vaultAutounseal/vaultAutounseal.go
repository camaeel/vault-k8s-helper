package main

import (
	"flag"

	"github.com/camaeel/vault-k8s-helper/pkg/autounseal"
)

func main() {
	cfg := autounseal.Config{}

	flag.StringVar(&cfg.ServiceDomain, "service-domain", "vault-internal.vault.svc.cluster.local", "DNS Name for accessing vault. In HA mode should be set to vault headles service providing all pod endpoints.")
	flag.StringVar(&cfg.ServiceScheme, "service-scheme", "https", "Vaul service scheme. Valid values: http, https")
	flag.IntVar(&cfg.ServicePort, "service-port", 8200, "Vaul service port")
	flag.IntVar(&cfg.UnlockShares, "unlock-shares", 3, "Number of unlock shares")
	flag.IntVar(&cfg.UnlockThreshold, "unlock-threshold", 2, "Number of unlock shares threshold")
	flag.StringVar(&cfg.VaultRootTokenSecret, "vault-root-token-secret", "vault-root-token", "Vault root token secret name")
	flag.StringVar(&cfg.VaultUnlockKeysSecret, "vault-unlock-keys-secret", "vault-unlock-keys", "Vault unlock keys secret name")

	flag.Parse()
	// ctx := context.TODO()

	for {
		autounseal.ManageVaultAutounseal(cfg)
	}

}
