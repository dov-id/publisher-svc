package indexer

import (
	"context"
	"fmt"
	"io"

	"github.com/dov-id/publisher-svc/contracts"
	"github.com/dov-id/publisher-svc/internal/config"
	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/dov-id/publisher-svc/internal/data/postgres"
	"github.com/dov-id/publisher-svc/internal/helpers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	shell "github.com/ipfs/go-ipfs-api"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const (
	serviceName = "indexer"
)

func NewIndexer(cfg config.Config) Indexer {
	return &indexer{
		cfg:              cfg,
		log:              cfg.Log(),
		FeedbacksQ:       postgres.NewFeedbacksQ(cfg.DB().Clone()),
		Shell:            shell.NewShell(cfg.Ipfs().Url),
		Clients:          map[string]*ethclient.Client{},
		FeedbackRegistry: map[string]*contracts.FeedbackRegistry{},
	}
}

func (i *indexer) Run(ctx context.Context) {
	go running.WithBackOff(
		ctx,
		i.log,
		serviceName,
		i.listen,
		i.cfg.Indexer().Timeout,
		i.cfg.Indexer().Timeout,
		i.cfg.Indexer().Timeout,
	)
}

func (i *indexer) listen(_ context.Context) error {
	i.log.Debugf("start feedback indexation")

	var err error
	i.Clients, err = helpers.InitNetworkClients(i.cfg.Networks().Networks)
	if err != nil {
		return errors.Wrap(err, "failed to init network clients")
	}

	i.FeedbackRegistry, err = helpers.InitFeedbackRegistryContracts(i.cfg.FeedbackRegistry().Addresses, i.Clients)
	if err != nil {
		return errors.Wrap(err, "failed to init feedback registry contracts")
	}

	//err = i.getFeedbacks()
	if err != nil {
		return errors.Wrap(err, "failed to get feedbacks")
	}

	i.log.Debugf("finish feedback indexation")
	return nil
}

func (i *indexer) getFeedbacks() error {
	for network, contract := range i.FeedbackRegistry {
		err := i.processGetFeedbacks(contract)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to get feebacks from `%s`", network))
		}
	}

	return nil
}

func (i *indexer) processGetFeedbacks(feedbackRegistry *contracts.FeedbackRegistry) error {
	response, err := feedbackRegistry.GetAllFeedbacks(new(bind.CallOpts))
	if err != nil {
		return errors.Wrap(err, "failed to get all feedbacks")
	}

	for k, course := range response.Courses {
		for _, feedbackBytes32 := range response.Feedbacks[k] {
			feedbackString, err := i.readFeedbackFromIPFS(common.Bytes2Hex(feedbackBytes32[:]))
			if err != nil {
				return errors.Wrap(err, "failed to read feedback from ipfs")
			}

			if feedbackString == nil {
				return errors.New(data.EmptyFeedbackContentErr)
			}

			err = i.FeedbacksQ.Insert(data.Feedback{
				Course:  common.Bytes2Hex(course),
				Content: *feedbackString,
			})
			if err != nil {
				return errors.Wrap(err, "failed to insert feedback")
			}
		}
	}

	return nil
}

func (i *indexer) readFeedbackFromIPFS(ipfsHash string) (*string, error) {
	reader, err := i.Shell.Cat(fmt.Sprintf("/ipfs/%s", ipfsHash))
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
