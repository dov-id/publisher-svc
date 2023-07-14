package handlers

import (
	"context"
	"net/http"

	"github.com/dov-id/publisher-svc/internal/config"
	"github.com/dov-id/publisher-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	cfgCtxKey
	requestsCtxKey
	feedbacksCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxRequestsQ(entry data.Requests) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, requestsCtxKey, entry)
	}
}

func RequestsQ(r *http.Request) data.Requests {
	return r.Context().Value(requestsCtxKey).(data.Requests)
}

func CtxFeedbacksQ(entry data.Feedbacks) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, feedbacksCtxKey, entry)
	}
}

func FeedbacksQ(r *http.Request) data.Feedbacks {
	return r.Context().Value(feedbacksCtxKey).(data.Feedbacks)
}

func CtxCfg(entry config.Config) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, cfgCtxKey, entry)
	}
}

func Cfg(r *http.Request) config.Config {
	return r.Context().Value(cfgCtxKey).(config.Config)
}
