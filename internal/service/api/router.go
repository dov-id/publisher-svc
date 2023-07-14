package api

import (
	"github.com/dov-id/publisher-svc/internal/data/postgres"
	"github.com/dov-id/publisher-svc/internal/service/api/handlers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *Router) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxCfg(s.cfg),
			handlers.CtxRequestsQ(postgres.NewRequestsQ(s.cfg.DB().Clone())),
			handlers.CtxFeedbacksQ(postgres.NewFeedbacksQ(s.cfg.DB().Clone())),
		),
	)
	r.Route("/integrations/publisher-svc", func(r chi.Router) {
		r.Post("/ring", handlers.GenerateRingSignature)

		r.Post("/requests", handlers.GetRequest)

		r.Route("/feedbacks", func(r chi.Router) {
			r.Post("/", handlers.AddFeedback)
			r.Get("/", handlers.GetFeedbacks)
		})
	})

	return r
}
