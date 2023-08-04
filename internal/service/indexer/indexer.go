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
	"github.com/ethereum/go-ethereum/common"
	ipfs "github.com/ipfs/go-ipfs-api"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const (
	serviceName = "indexer"
)

func NewIndexer(cfg config.Config) Indexer {
	return &indexer{
		cfg:                 cfg,
		log:                 cfg.Log(),
		FeedbacksQ:          postgres.NewFeedbacksQ(cfg.DB().Clone()),
		Ipfs:                ipfs.NewShell(cfg.Ipfs().Url),
		LastHandledFeedback: make(map[common.Address]int64),
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

	err := i.startFeedbackIndexation(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get feedbacks")
	}

	i.log.Debugf("finish feedback indexation")
	return nil
}

func (i *indexer) startFeedbackIndexation(ctx context.Context) error {
	feedbackRegistries, err := helpers.GetFeedbackRegistriesFromCtx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get feedbacks registries from ctx")
	}

	for network, contract := range feedbackRegistries {
		err = i.indexFeedbacks(contract)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to get feebacks from `%s`", network))
		}
	}

	return nil
}

func (i *indexer) indexFeedbacks(feedbackRegistry *contracts.FeedbackRegistry) error {
	err := i.getCourses(feedbackRegistry)
	if err != nil {
		return errors.Wrap(err, "failed to get courses")
	}

	err = i.getFeedbacks(feedbackRegistry)
	if err != nil {
		return errors.Wrap(err, "failed to get feedbacks")
	}

	return nil
}

func (i *indexer) getCourses(feedbackRegistry *contracts.FeedbackRegistry) error {
	var offset = big.NewInt(0)
	var limit = big.NewInt(15)

	for {
		courses, err := feedbackRegistry.GetCourses(new(bind.CallOpts), offset, limit)
		if err != nil {
			return errors.Wrap(err, "failed to get courses")
		}

		if len(courses) == 0 {
			break
		}

		for _, course := range courses {
			_, isPresent := i.LastHandledFeedback[course]
			if !isPresent {
				i.LastHandledFeedback[course] = 0
			}
		}

		offset.Add(offset, limit)
	}

	return nil
}

func (i *indexer) getFeedbacks(feedbackRegistry *contracts.FeedbackRegistry) error {
	for course, lastFeedbackNumber := range i.LastHandledFeedback {
		var offset = big.NewInt(lastFeedbackNumber)
		var limit = big.NewInt(15)

		for {
			feedbacks, err := feedbackRegistry.GetFeedbacks(new(bind.CallOpts), course, offset, limit)
			if err != nil {
				return errors.Wrap(err, "failed to get all feedbacks")
			}

			if len(feedbacks) == 0 {
				break
			}

			for _, feedback := range feedbacks {
				err = i.processFeedback(course, feedback)
				if err != nil {
					return errors.Wrap(err, "failed to process feedback")
				}
			}

			offset.Add(offset, limit)
		}
	}

	return nil
}

func (i *indexer) processFeedback(course common.Address, feedback string) error {
	feedbackString, err := i.readFeedbackFromIPFS(feedback)
	if err != nil {
		return errors.Wrap(err, "failed to read feedback from ipfs")
	}

	if feedbackString == nil {
		return data.ErrEmptyFeedbackContent
	}

	err = i.FeedbacksQ.Insert(data.Feedback{
		Course:  course.Hex(),
		Content: *feedbackString,
	})
	if err != nil {
		return errors.Wrap(err, "failed to insert feedback")
	}

	i.LastHandledFeedback[course]++

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
