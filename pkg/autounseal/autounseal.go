package autounseal

import (
	"context"
	"fmt"

	"github.com/camaeel/vault-k8s-helper/pkg/config"
	"github.com/camaeel/vault-k8s-helper/pkg/k8sProvider"
	"github.com/camaeel/vault-k8s-helper/pkg/providers/vaultClient"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"k8s.io/client-go/kubernetes"
)

func ManageVaultAutounseal(cfg config.Config, ctx context.Context, k8s *kubernetes.Clientset) error {

	nodeNames, err := k8sProvider.GetServiceEndpoints(k8s, ctx, &cfg.VaultInternalServiceName, &cfg.VaultNamespace)
	if err != nil {
		return err
	}

	expectedFirstNodeName := fmt.Sprintf("%s-0", cfg.VaultPodNamePrefix)
	if nodeNames[0] != expectedFirstNodeName {
		return fmt.Errorf("First node %s is not equal to expected %s", nodeNames[0], expectedFirstNodeName)
	}

	nodes := make([]*vaultClient.Node, 0)
	for n := range nodeNames {
		node, err := vaultClient.GetVaultClusterNode(ctx, podEndpoint(cfg, nodeNames[n]), cfg)
		if err != nil {
			return err
		}
		nodes = append(nodes, node)
	}

	for n := range nodes {

		if !nodes[n].Initialized {
			if n == 0 {
				//first node & needs initialization
				log.Info("Initializing node %d", n)
				keys, rootToken, err := nodes[n].Initialize(cfg, ctx)
				if err != nil {
					return err
				}
				// join for others - how?
				// https://developer.hashicorp.com/vault/docs/platform/k8s/helm/examples/ha-with-raft

				keysMap := mapVaultKeys(keys)
				rootTokenMap := map[string][]byte{
					"token": []byte(rootToken),
				}

				log.Infof("creating secrets containg initialziation data")
				k8sProvider.CreateOrReplaceSecret(k8s, ctx, &cfg.VaultUnlockKeysSecret, &cfg.Namespace, keysMap)
				k8sProvider.CreateOrReplaceSecret(k8s, ctx, &cfg.VaultRootTokenSecret, &cfg.Namespace, rootTokenMap)
				log.Infof("secrets containg initialziation data created")
			} else {
				log.Infof("Joining node %d to existing cluster", n)
				err := nodes[n].Join(cfg, ctx, nodes[0])
				if err != nil {
					return err
				}
				log.Infof("Joined node %d to existing cluster successfully", n)
			}
		} else {
			log.Infof("Node %d already initialized.", n)
			if nodes[n].Sealed {
				log.Infof("Node %d sealed. Will try to unseal", n)
				keysMap, err := k8sProvider.GetSecretContents(k8s, ctx, &cfg.VaultUnlockKeysSecret, &cfg.Namespace)
				if err != nil {
					return nil
				}
				keys := maps.Values(keysMap)
				for k := range keys {
					nodes[n].Unseal(ctx, string(keys[k]), k)
				}

			} else {
				log.Infof("Node %d initialized and unsealed. Nothing to do", n)
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

func podEndpoint(cfg config.Config, nodeName string) string {
	return fmt.Sprintf("%s://%s.%s:%d", cfg.ServiceScheme, nodeName, cfg.ServiceDomain, cfg.ServicePort)
}
