package helpers

import (
	"fmt"

	"github.com/dov-id/publisher-svc/contracts"
	"github.com/dov-id/publisher-svc/internal/config"
	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func InitNetworkClients(networks map[string]config.Network) (map[string]*ethclient.Client, error) {
	clients := make(map[string]*ethclient.Client)

	infura := networks[data.InfuraNetwork]

	for network, params := range networks {
		if network == data.InfuraNetwork || network == data.MetamaskNetwork {
			continue
		}

		client, err := CreateNetworkClient(network, params, infura)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create network client")
		}

		clients[network] = client
	}

	return clients, nil
}

func CreateNetworkClient(network string, params config.Network, infura config.Network) (*ethclient.Client, error) {
	var rawUrl string
	switch network {
	case data.QNetwork:
		rawUrl = params.RpcUrl
	default:
		rawUrl = params.RpcUrl + infura.Key
	}

	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to make dial connect to `%s` network", network))
	}

	return client, nil
}

func InitFeedbackRegistryContracts(certIntegrators map[string]string, clients map[string]*ethclient.Client) (map[string]*contracts.FeedbackRegistry, error) {
	certIntegratorContracts := make(map[string]*contracts.FeedbackRegistry)

	for network, address := range certIntegrators {
		contract, err := contracts.NewFeedbackRegistry(common.HexToAddress(address), clients[network])
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to create new `%s` feedback registry contract", network))
		}

		certIntegratorContracts[network] = contract
	}

	return certIntegratorContracts, nil
}
