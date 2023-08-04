package helpers

import (
	"context"

	"github.com/dov-id/publisher-svc/contracts"
	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetNetworkClientsFromCtx(ctx context.Context) (map[data.Network]*ethclient.Client, error) {
	value := ctx.Value(data.NetworkClients)

	clients, ok := value.(map[data.Network]*ethclient.Client)
	if !ok {
		return nil, data.ErrFailedToCastClients
	}

	return clients, nil
}

func GetFeedbackRegistriesFromCtx(ctx context.Context) (map[data.Network]*contracts.FeedbackRegistry, error) {
	value := ctx.Value(data.FeedbackRegistriesContracts)

	feedbackRegistries, ok := value.(map[data.Network]*contracts.FeedbackRegistry)
	if !ok {
		return nil, data.ErrFailedToCastFeedbackRegistries
	}

	return feedbackRegistries, nil
}
