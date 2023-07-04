package helpers

import (
	"context"
	"math/big"

	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func GetAuth(client *ethclient.Client, private string) (*bind.TransactOpts, error) {
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chain id")
	}

	privateKey, fromAddress, err := GetKeys(private)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get keys")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transaction signer")
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get nonce")
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to suggest gas price")
	}

	auth.GasLimit = uint64(10000000)
	auth.GasPrice = gasPrice

	auth.Nonce = big.NewInt(int64(nonce))

	return auth, nil
}

func WaitForTransactionMined(client *ethclient.Client, transaction *types.Transaction, log *logan.Entry, reqId string, requestsQ data.Requests) {
	var (
		err   error
		mined = make(chan struct{})
		ctx   = context.Background()
	)

	go func() {
		log.WithField("tx", transaction.Hash().Hex()).Debugf("waiting to mine")

		_, err = bind.WaitMined(ctx, client, transaction)
		if err != nil {
			errorMsg := err.Error()

			err = requestsQ.FilterByIds(reqId).Update(data.RequestToUpdate{Status: data.FAILED, Error: &errorMsg})
			if err != nil {
				err = errors.Wrap(err, "failed to update request status")
			}

			panic(errors.Wrap(err, "failed to mine transaction"))
		}

		err = requestsQ.FilterByIds(reqId).Update(data.RequestToUpdate{Status: data.SUCCESS})
		if err != nil {
			panic(errors.Wrap(err, "failed to update request status"))
		}

		log.WithField("tx", transaction.Hash().Hex()).Debugf("was mined")

		close(mined)
	}()
}
