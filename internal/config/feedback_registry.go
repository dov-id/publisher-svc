package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type FeedbackRegistryCfg struct {
	Addresses map[string]string
}

type feedbacksCfg struct {
	List []feedback
}

type feedback struct {
	Network string `fig:"network,required"`
	Address string `fig:"address,required"`
}

func (c *config) FeedbackRegistry() *FeedbackRegistryCfg {
	return c.feedbackRegistry.Do(func() interface{} {
		var cfg feedbacksCfg

		err := figure.
			Out(&cfg).
			With(figure.BaseHooks, FeedbackHooks).
			From(kv.MustGetStringMap(c.getter, "feedback_registry")).
			Please()

		if err != nil {
			panic(errors.Wrap(err, "failed to figure out feedback registry config"))
		}

		mapCfg := createMapIntegrators(cfg.List)
		return &mapCfg
	}).(*FeedbackRegistryCfg)
}

func createMapIntegrators(list []feedback) FeedbackRegistryCfg {
	var cfg FeedbackRegistryCfg
	cfg.Addresses = make(map[string]string)

	for _, elem := range list {
		cfg.Addresses[elem.Network] = elem.Address
	}

	return cfg
}
