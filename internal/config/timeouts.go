package config

import (
	"time"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type TimeoutsCfg struct {
	Indexer time.Duration `figure:"indexer,required"`
}

func (c *config) Timeouts() *TimeoutsCfg {
	return c.timeouts.Do(func() interface{} {
		var cfg TimeoutsCfg

		err := figure.
			Out(&cfg).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "timeouts")).
			Please()

		if err != nil {
			panic(errors.Wrap(err, "failed to figure out timeouts config"))
		}

		return &cfg
	}).(*TimeoutsCfg)
}
