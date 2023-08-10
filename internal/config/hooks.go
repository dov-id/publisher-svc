package config

import (
	"fmt"
	"reflect"

	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var NetworkHooks = figure.Hooks{
	"config.NetworksCfg": func(value interface{}) (reflect.Value, error) {
		type network struct {
			Name                    string `fig:"name,required"`
			RpcProviderWsUrl        string `fig:"rpc_provider_ws_url,required"`
			FeedbackRegistryAddress string `fig:"feedback_registry_address,required"`
			WalletPrivateKey        string `fig:"wallet_private_key,required"`
		}

		switch v := value.(type) {
		case map[string]interface{}:
			cfg := NetworksCfg{}
			cfg.Networks = make(map[data.Network]Network)

			for _, rawMap := range v {
				parsedMapping, ok := rawMap.([]interface{})
				if !ok {
					return reflect.Value{}, errors.New("failed to cast raw map element to []interface{}")
				}

				for _, mapElement := range parsedMapping {
					mapElem, ok := mapElement.(map[interface{}]interface{})
					if !ok {
						return reflect.Value{}, errors.New("failed to cast map element to map[interface{}]interface{}")

					}

					normMap := make(map[string]interface{}, len(mapElem))

					for key, value := range mapElem {
						normMap[fmt.Sprint(key)] = value
					}

					var info network

					err := figure.
						Out(&info).
						With(figure.BaseHooks).
						From(normMap).
						Please()
					if err != nil {
						return reflect.Value{}, errors.Wrap(err, "failed to figure out network data")
					}

					if !common.IsHexAddress(info.FeedbackRegistryAddress) {
						return reflect.Value{}, errors.Wrap(err, "cert integrator hex is not an address")
					}

					privateKey, fromAddress, err := getKeys(info.WalletPrivateKey)
					if err != nil {
						panic(errors.Wrap(err, "failed to retrieve keys"))
					}

					cfg.Networks[data.Network(info.Name)] = Network{
						RpcProviderWsUrl:        info.RpcProviderWsUrl,
						FeedbackRegistryAddress: common.HexToAddress(info.FeedbackRegistryAddress),
						WalletCfg: &WalletCfg{
							PrivateKey: privateKey,
							Address:    fromAddress,
						},
					}
				}
			}

			return reflect.ValueOf(cfg), nil
		default:
			return reflect.Value{}, fmt.Errorf("unexpected type to figure Config.Network[]")
		}
	},
}
