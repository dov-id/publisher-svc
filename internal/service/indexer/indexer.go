package indexer

import (
	"context"
	"fmt"
	"io"
	"math/big"

	"github.com/dov-id/publisher-svc/contracts"
	"github.com/dov-id/publisher-svc/internal/config"
	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/dov-id/publisher-svc/internal/data/postgres"
	"github.com/dov-id/publisher-svc/internal/helpers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ipfs "github.com/ipfs/go-ipfs-api"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const (
	serviceName = "indexer"
)

func NewIndexer(cfg config.Config) Indexer {
	return &indexer{
		cfg:                   cfg,
		log:                   cfg.Log(),
		FeedbacksQ:            postgres.NewFeedbacksQ(cfg.DB().Clone()),
		Ipfs:                  ipfs.NewShell(cfg.Ipfs().Url),
		LastHandledFeedbackId: 0,
	}
}

func (i *indexer) Run(ctx context.Context) {
	go running.WithBackOff(
		ctx,
		i.log,
		serviceName,
		i.listen,
		i.cfg.Timeouts().Indexer,
		i.cfg.Timeouts().Indexer,
		i.cfg.Timeouts().Indexer,
	)
}

func (i *indexer) listen(ctx context.Context) error {
	i.log.Debugf("start feedback indexation")

	err := i.getFeedbacks(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get feedbacks")
	}

	i.log.Debugf("finish feedback indexation")
	return nil
}

func (i *indexer) getFeedbacks(ctx context.Context) error {
	feedbackRegistries, err := helpers.GetFeedbackRegistriesFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get feedbacks registries from ctx")
	}

	for network, contract := range feedbackRegistries {
		err = i.processGetFeedbacks(contract)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to get feebacks from `%s`", network))
		}
	}

	return nil
}

func (i *indexer) processGetFeedbacks(feedbackRegistry *contracts.FeedbackRegistry) error {
	var offset = big.NewInt(i.LastHandledFeedbackId)
	var limit = big.NewInt(15)

	for {
		response, err := feedbackRegistry.GetAllFeedbacks(new(bind.CallOpts), offset, limit)
		if err != nil {
			return errors.Wrap(err, "failed to get all feedbacks")
		}

		if len(response.Courses) == 0 {
			break
		}

		for k, course := range response.Courses {
			for _, feedback := range response.Feedbacks[k] {
				err = i.processFeedback(int64(k), course.Hex(), feedback)
				if err != nil {
					return errors.Wrap(err, "failed to process feedback")
				}
			}
		}

		offset.Add(offset, limit)
	}

	return nil
}

func (i *indexer) processFeedback(index int64, course, feedback string) error {
	feedbackString, err := i.readFeedbackFromIPFS(feedback)
	if err != nil {
		return errors.Wrap(err, "failed to read feedback from ipfs")
	}

	if feedbackString == nil {
		return data.ErrEmptyFeedbackContent
	}

	err = i.FeedbacksQ.Insert(data.Feedback{
		Course:  course,
		Content: *feedbackString,
	})
	if err != nil {
		return errors.Wrap(err, "failed to insert feedback")
	}

	i.LastHandledFeedbackId = index

	return nil
}

func (i *indexer) readFeedbackFromIPFS(ipfsHash string) (*string, error) {
	reader, err := i.Ipfs.Cat(fmt.Sprintf("/ipfs/%s", ipfsHash))
	if err != nil {
		return nil, errors.Wrap(err, "failed to cat feedback by cid")
	}

	bytesContent, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read bytes feedback content")
	}

	feedback := string(bytesContent)

	return &feedback, nil
}
