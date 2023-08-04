package data

import (
	"errors"
)

type Network string

func (n Network) String() string { return string(n) }

var (
	ErrFailedToSetString              = errors.New("failed to set big int string")
	ErrFailedToCastKey                = errors.New("failed to cast public key to ECDSA")
	ErrFailedToCastFeedbackRegistries = errors.New("failed to cast public interface{} to map[types.Network]*contracts.FeedbackRegistry")
	ErrFailedToCastClients            = errors.New("failed to cast public interface{} to map[types.Network]*ethclient.Client")
	ErrEmptyFeedbackContent           = errors.New("feedback content is empty")
	ErrReplacementTxUnderpriced       = errors.New("replacement transaction underpriced")
	ErrNotString                      = errors.New("the value is not a string")
	ErrInvalidEthAddress              = errors.New("given value is invalid ethereum address")
)

const (
	EthereumNetwork Network = "ethereum"
	PolygonNetwork  Network = "polygon"
	QNetwork        Network = "q"
)

const (
	NetworkClients              = "network clients"
	FeedbackRegistriesContracts = "feedback registries contracts"
)
