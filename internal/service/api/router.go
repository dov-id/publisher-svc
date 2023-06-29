package api

import (
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
		),
	)
	r.Route("/integrations/publisher-svc", func(r chi.Router) {
		r.Post("/ring", handlers.GenerateRingSignature)

		r.Route("/feedbacks", func(r chi.Router) {
			r.Post("/", handlers.AddFeedback)
		})
	})

	return r
}
