package main

import (
	"context"
	"flag"
	"time"

	"github.com/camaeel/vault-k8s-helper/pkg/autounseal"
	"github.com/camaeel/vault-k8s-helper/pkg/config"
	"github.com/camaeel/vault-k8s-helper/pkg/k8sProvider"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Config{}

	flag.StringVar(&cfg.ServiceDomain, "service-domain", "vault-internal.vault.svc.cluster.local", "DNS Name for accessing vault. In HA mode should be set to vault headles service providing all pod endpoints.")
	flag.StringVar(&cfg.ServiceScheme, "service-scheme", "https", "Vaul service scheme. Valid values: http, https")
	flag.IntVar(&cfg.ServicePort, "service-port", 8200, "Vaul service port")
	flag.IntVar(&cfg.UnlockShares, "unlock-shares", 3, "Number of unlock shares")
	flag.IntVar(&cfg.UnlockThreshold, "unlock-threshold", 3, "Number of unlock shares threshold")
	flag.StringVar(&cfg.VaultRootTokenSecret, "vault-root-token-secret", "vault-autounseal-root-token", "Vault root token secret name")
	flag.StringVar(&cfg.VaultUnlockKeysSecret, "vault-unlock-keys-secret", "vault-autounseal-unlock-keys", "Vault unlock keys secret name")
	flag.StringVar(&cfg.Namespace, "namespace", "vault-autounseal", "Namespace used for storing unseal keys and root token")
	flag.StringVar(&cfg.CaCert, "ca-cert", "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt", "CA certificate for validating connections to vault")
	flag.StringVar(&cfg.VaultPodNamePrefix, "vault-pod-name-prefix", "vault", "Prefix for vault StatefulSet's pods")
	flag.StringVar(&cfg.VaultInternalServiceName, "vault-internal-service-name", "vault-internal", "Name of vault's internal service name")
	flag.StringVar(&cfg.VaultNamespace, "vault-namespace", "vault", "namespace where vault is installed")
	kubeconfig := flag.String("kubeconfig", "", "Overwrite kubeconfig path")

	flag.Parse()
	ctx := context.TODO()

	k8s, err := k8sProvider.GetClientSet(kubeconfig)
	if err != nil {
		panic(err)
	}

	for {
		err := autounseal.ManageVaultAutounseal(cfg, ctx, k8s)
		if err != nil {
			log.Errorf("received error: %v", err)
		}
		log.Debug("Sleeping")
		time.Sleep(30 * time.Second)
	}

}
