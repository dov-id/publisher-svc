package indexer

import (
	"context"

	"github.com/dov-id/publisher-svc/internal/config"
	"github.com/dov-id/publisher-svc/internal/data"
	ipfs "github.com/ipfs/go-ipfs-api"
	"gitlab.com/distributed_lab/logan/v3"
)

type Indexer interface {
	Run(ctx context.Context)
}

type indexer struct {
	cfg config.Config
	log *logan.Entry

	LastHandledFeedbackId int64
	FeedbacksQ            data.Feedbacks
	Ipfs                  *ipfs.Shell
}

func Run(cfg config.Config, ctx context.Context) {
	NewIndexer(cfg).Run(ctx)
}
