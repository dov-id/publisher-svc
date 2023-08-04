package service

import (
	"context"
	"sync"

	"github.com/dov-id/publisher-svc/internal/config"
	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/dov-id/publisher-svc/internal/helpers"
	"github.com/dov-id/publisher-svc/internal/service/api"
	"github.com/dov-id/publisher-svc/internal/service/indexer"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Runner = func(config config.Config, context context.Context)

var availableServices = map[string]Runner{
	"api":     api.Run,
	"indexer": indexer.Run,
}

func Run(ctx context.Context, cfg config.Config) {
	logger := cfg.Log().WithField("service", "main")

	ctx, err := prepareContextStorage(ctx, cfg)
	if err != nil {
		panic(errors.Wrap(err, "failed to prepare service context storage"))
	}

	wg := new(sync.WaitGroup)

	logger.Debugf("Starting all available services")

	for serviceName, service := range availableServices {
		wg.Add(1)

		go func(name string, runner Runner) {
			defer wg.Done()

			runner(cfg, ctx)

		}(serviceName, service)

		logger.WithField("service", serviceName).Debugf("Service started")
	}

	wg.Wait()
}

func prepareContextStorage(ctx context.Context, cfg config.Config) (context.Context, error) {
	clients, err := helpers.InitNetworkClients(cfg.Networks().Networks)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize network clients")
	}

	ctx = context.WithValue(ctx, data.NetworkClients, clients)

	feedbackRegistries, err := helpers.InitFeedbackRegistryContracts(cfg.Networks().Networks, clients)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize feedback registry contracts")
	}

	ctx = context.WithValue(ctx, data.FeedbackRegistriesContracts, feedbackRegistries)

	return ctx, nil
}
