// Copyright 2025 Erst Users
// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
)

// NetworkConfig represents a Stellar network configuration
type NetworkConfig struct {
	Name              string
	HorizonURL        string
	NetworkPassphrase string
	SorobanRPCURL     string
}

// Predefined network configurations
var (
	TestnetConfig = NetworkConfig{
		Name:              "testnet",
		HorizonURL:        "https://horizon-testnet.stellar.org/",
		NetworkPassphrase: "Test SDF Network ; September 2015",
		SorobanRPCURL:     "https://soroban-testnet.stellar.org",
	}

	MainnetConfig = NetworkConfig{
		Name:              "mainnet",
		HorizonURL:        "https://horizon.stellar.org/",
		NetworkPassphrase: "Public Global Stellar Network ; September 2015",
		SorobanRPCURL:     "https://mainnet.stellar.validationcloud.io/v1/soroban-rpc-demo",
	}

	FuturenetConfig = NetworkConfig{
		Name:              "futurenet",
		HorizonURL:        "https://horizon-futurenet.stellar.org/",
		NetworkPassphrase: "Test SDF Future Network ; October 2022",
		SorobanRPCURL:     "https://rpc-futurenet.stellar.org",
	}
)

// Client handles interactions with the Stellar Network
type Client struct {
	Horizon *horizonclient.Client
	Config  NetworkConfig
}

// NewClient creates a new RPC client for a predefined network
func NewClient(networkName string) (*Client, error) {
	var config NetworkConfig

	switch networkName {
	case "testnet":
		config = TestnetConfig
	case "mainnet", "public":
		config = MainnetConfig
	case "futurenet":
		config = FuturenetConfig
	default:
		return nil, fmt.Errorf("unknown network: %s (use 'testnet', 'mainnet', or 'futurenet')", networkName)
	}

	return &Client{
		Horizon: &horizonclient.Client{
			HorizonURL: config.HorizonURL,
			HTTP:       http.DefaultClient,
		},
		Config: config,
	}, nil
}

// NewCustomClient creates a new RPC client for a custom/private network
func NewCustomClient(config NetworkConfig) (*Client, error) {
	if config.HorizonURL == "" {
		return nil, fmt.Errorf("horizon URL is required for custom network")
	}
	if config.NetworkPassphrase == "" {
		return nil, fmt.Errorf("network passphrase is required for custom network")
	}

	return &Client{
		Horizon: &horizonclient.Client{
			HorizonURL: config.HorizonURL,
			HTTP:       http.DefaultClient,
		},
		Config: config,
	}, nil
}

// TransactionResponse contains the raw XDR fields needed for simulation
type TransactionResponse struct {
	EnvelopeXdr   string
	ResultXdr     string
	ResultMetaXdr string
}

// GetTransaction fetches the transaction details and full XDR data
func (c *Client) GetTransaction(ctx context.Context, hash string) (*TransactionResponse, error) {
	tx, err := c.Horizon.TransactionDetail(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}

	return &TransactionResponse{
		EnvelopeXdr:   tx.EnvelopeXdr,
		ResultXdr:     tx.ResultXdr,
		ResultMetaXdr: tx.ResultMetaXdr,
	}, nil
}

// GetNetworkPassphrase returns the network passphrase for this client
func (c *Client) GetNetworkPassphrase() string {
	return c.Config.NetworkPassphrase
}

// GetNetworkName returns the network name for this client
func (c *Client) GetNetworkName() string {
	if c.Config.Name != "" {
		return c.Config.Name
	}
	return "custom"
}
