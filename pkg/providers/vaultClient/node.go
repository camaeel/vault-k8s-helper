package vaultClient

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/camaeel/vault-k8s-helper/pkg/config"
	vault "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
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

func (n *Node) Join(cfg config.Config, ctx context.Context, node0 *Node) error {
	cacert, err := ioutil.ReadFile(cfg.CaCert)
	if err != nil {
		return err
	}
	input := vault.RaftJoinRequest{
		LeaderAPIAddr: node0.Address,
		LeaderCACert:  string(cacert),
	}
	resp, err := n.Client.Sys().RaftJoinWithContext(ctx, &input)
	if err != nil {
		return err
	}
	if !resp.Joined {
		return fmt.Errorf("Raft join of node %s to %s not successful", n.Address, node0.Address)
	}
	return nil
}
