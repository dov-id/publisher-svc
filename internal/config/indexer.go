package config

import (
	"time"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type IndexerCfg struct {
	Timeout time.Duration `figure:"timeout,required"`
}

func (c *config) Indexer() *IndexerCfg {
	return c.indexer.Do(func() interface{} {
		var cfg IndexerCfg

		err := figure.
			Out(&cfg).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "indexer")).
			Please()

		if err != nil {
			panic(errors.Wrap(err, "failed to figure out indexer config"))
		}

		return &cfg
	}).(*IndexerCfg)
}
