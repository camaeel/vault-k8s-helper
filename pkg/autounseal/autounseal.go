package autounseal

import (
	"context"
	"fmt"
	"net"

	"github.com/camaeel/vault-k8s-helper/pkg/config"
	"github.com/camaeel/vault-k8s-helper/pkg/k8sProvider"
	"github.com/camaeel/vault-k8s-helper/pkg/providers/vaultClient"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"k8s.io/client-go/kubernetes"
)

func ManageVaultAutounseal(cfg config.Config, ctx context.Context, k8s *kubernetes.Clientset) error {
	ips, err := net.LookupIP(cfg.ServiceDomain)
	if err != nil {
		return err
	}
	log.Infof("Found ips: %v", ips)

	nodes := make([]*vaultClient.Node, 0)
	for i := range ips {
		address := fmt.Sprintf("%s://%s:%d", cfg.ServiceScheme, ips[i].String(), cfg.ServicePort)
		node, err := vaultClient.GetVaultClusterNode(ctx, address)
		if err != nil {
			return err
		}
		nodes = append(nodes, node)
	}

	atLeastOneInitialized := false
	allInitialized := true

	for i := range nodes {
		atLeastOneInitialized = atLeastOneInitialized || nodes[i].Initialized
		allInitialized = allInitialized && nodes[i].Initialized
	}

	if !atLeastOneInitialized {
		// do initialize
		log.Info("no nodes initialized")
		keys, rootToken, err := nodes[0].Initialize(cfg, ctx)
		if err != nil {
			return err
		}

		keysMap := mapVaultKeys(keys)
		rootTokenMap := map[string][]byte{
			"token": []byte(rootToken),
		}

		log.Info("creating secrets containg initialziation data")
		k8sProvider.CreateOrReplaceSecret(k8s, ctx, &cfg.VaultUnlockKeysSecret, &cfg.Namespace, keysMap)
		k8sProvider.CreateOrReplaceSecret(k8s, ctx, &cfg.VaultRootTokenSecret, &cfg.Namespace, rootTokenMap)

	} else if !allInitialized { //only some are initialzied - shouldn't this be an error ???
		log.Warn("Only some nodes are initialized")
	} else {
		log.Info("All nodes initialized")
		keysMap, err := k8sProvider.GetSecretContents(k8s, ctx, &cfg.VaultUnlockKeysSecret, &cfg.Namespace)
		if err != nil {
			return err
		}
		for i := range nodes {
			if !nodes[i].Sealed {
				keys := maps.Values(keysMap)
				for k := range keys {
					nodes[i].Unseal(ctx, string(keys[k]), k)
				}
			}
		}
	}

	return nil
}

func mapVaultKeys(keys []string) map[string][]byte {
	ret := make(map[string][]byte)
	for k := range keys {
		ret[fmt.Sprintf("key-%d", k)] = []byte(keys[k])
	}
	return ret
}
