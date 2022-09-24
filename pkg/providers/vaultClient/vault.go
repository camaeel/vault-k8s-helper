package vaultClient

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

func newVaultClient(ctx context.Context, address string) (*vault.Client, error) {
	log.Debugf("connecting to vault @ %s", address)

	config := vault.DefaultConfig() // modify for more granular configuration
	config.Address = address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize vault client: %w", err)
	}

	return client, nil
}

func GetVaultClusterNode(ctx context.Context, address string) (*Node, error) {
	var err error

	node := Node{}
	node.Address = address
	node.Client, err = newVaultClient(ctx, node.Address)
	if err != nil {
		return nil, err
	}

	sealStatusResponse, err := node.Client.Sys().SealStatusWithContext(ctx)
	if err != nil {
		return nil, err
	}

	node.Sealed = sealStatusResponse.Sealed
	node.Initialized = sealStatusResponse.Initialized

	return &node, nil
}
