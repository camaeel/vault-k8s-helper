package vaultClient

import (
	"context"
	"fmt"

	"github.com/camaeel/vault-k8s-helper/pkg/config"
	vault "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

func newVaultClient(ctx context.Context, address string, cfg config.Config) (*vault.Client, error) {
	log.Debugf("connecting to vault @ %s", address)

	config := vault.DefaultConfig() // modify for more granular configuration
	config.Address = address

	tlsConfig := vault.TLSConfig{
		CACert: cfg.CaCert,
	}
	err := config.ConfigureTLS(&tlsConfig)
	if err != nil {
		return nil, err
	}

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize vault client: %w", err)
	}

	return client, nil
}

func GetVaultClusterNode(ctx context.Context, address string, cfg config.Config) (*Node, error) {
	var err error

	node := Node{}
	node.Address = address
	node.Client, err = newVaultClient(ctx, node.Address, cfg)
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
