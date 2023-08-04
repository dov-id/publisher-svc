package config

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type IpfsCfg struct {
	Url string `figure:"url,required"`
}

func (c *config) Ipfs() *IpfsCfg {
	return c.ipfs.Do(func() interface{} {
		var cfg IpfsCfg

		err := figure.
			Out(&cfg).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "ipfs")).
			Please()

		if err != nil {
			panic(errors.Wrap(err, "failed to figure out ipfs config"))
		}

		return &cfg
	}).(*IpfsCfg)
}
