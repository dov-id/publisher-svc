package helpers

import (
	"crypto/ecdsa"
	"sync"

	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func GetKeys(private string) (*ecdsa.PrivateKey, common.Address, error) {
	var once sync.Once
	var privateKey *ecdsa.PrivateKey
	var fromAddress common.Address
	var err error

	once.Do(func() {
		privateKey, err = crypto.HexToECDSA(private)
		if err != nil {
			err = errors.Wrap(err, "failed to convert hex to ecdsa")
			return
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			err = errors.New(data.FailedToCastKeyErr)
			return
		}

		fromAddress = crypto.PubkeyToAddress(*publicKeyECDSA)
		return
	})

	return privateKey, fromAddress, nil
}
