package vaultClient

import (
	"context"

	"github.com/camaeel/vault-k8s-helper/pkg/config"
	vault "github.com/hashicorp/vault/api"
	"github.com/prometheus/common/log"
)

type Node struct {
	Address     string
	Sealed      bool
	Initialized bool
	Client      *vault.Client
}

func (n *Node) Initialize(cfg config.Config, ctx context.Context) ([]string, string, error) {
	req := vault.InitRequest{
		SecretShares:    cfg.UnlockShares,
		SecretThreshold: cfg.UnlockThreshold,
	}
	resp, err := n.Client.Sys().InitWithContext(ctx, &req)
	if err != nil {
		return nil, "", err
	}
	keys := resp.Keys
	rootToken := resp.RootToken
	return keys, rootToken, nil
}

func (n *Node) Unseal(ctx context.Context, key string, keyIndex int) error {
	log.Infof("Unsealing vault node %s with key number %d", n.Address, keyIndex)
	_, err := n.Client.Sys().UnsealWithContext(ctx, key)

	return err
}
