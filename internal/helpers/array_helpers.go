package helpers

import (
	"github.com/dov-id/publisher-svc/crypto_master/secp256k1/ecc_math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func ConvertHexKeysToECPoints(publicKeys []string) ([]ecc_math.ECPoint, error) {
	var ECPoints = make([]ecc_math.ECPoint, len(publicKeys))
	for i, key := range publicKeys {
		if key[:2] == "0x" {
			key = key[2:]
		}

		publicECDSA, err := crypto.UnmarshalPubkey(common.Hex2Bytes(key))
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert hex to ecdsa")
		}

		ECPoints[i] = ecc_math.ECPoint{X: publicECDSA.X, Y: publicECDSA.Y}
	}
	return ECPoints, nil
}
