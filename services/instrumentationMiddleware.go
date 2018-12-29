package services

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           OAuth2Service
}

// Token Implements OAuth2Service
func (mw instrumentingMiddleware) Token(s tokenRequest) (output tokenResponse, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "token", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.next.Token(s)
	return
}

// Validate Implements OAuth2Service
func (mw instrumentingMiddleware) Validate(s validationRequest) (output validationResponse, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "validate", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.next.Validate(s)
	return
}

// NewInstrumentingMiddleware returns instrumenting middleware
func NewInstrumentingMiddleware(svc OAuth2Service) OAuth2Service {
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "goauth",
		Subsystem: "oauth2",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "goauth",
		Subsystem: "oauth2",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	return instrumentingMiddleware{requestCount, requestLatency, svc}
}
