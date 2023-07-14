package indexer

import (
	"context"

	"github.com/dov-id/publisher-svc/contracts"
	"github.com/dov-id/publisher-svc/internal/config"
	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/ethereum/go-ethereum/ethclient"
	shell "github.com/ipfs/go-ipfs-api"
	"gitlab.com/distributed_lab/logan/v3"
)

type Indexer interface {
	Run(ctx context.Context)
}

type indexer struct {
	cfg config.Config
	log *logan.Entry

	FeedbacksQ data.Feedbacks

	Shell            *shell.Shell
	Clients          map[string]*ethclient.Client
	FeedbackRegistry map[string]*contracts.FeedbackRegistry
}

func Run(cfg config.Config, ctx context.Context) {
	NewIndexer(cfg).Run(ctx)
}
