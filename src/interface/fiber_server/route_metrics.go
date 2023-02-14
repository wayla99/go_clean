package fiber_server

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const metricsPath = "/metrics"

type requestLabel struct {
	statusCode int
	method     string
	path       string
	duration   time.Duration
}

func (f *FiberServer) addMetrics(appVersion string) {
	hs, _ := os.Hostname()

	labels := prometheus.Labels{
		"version":  appVersion,
		"hostname": hs,
	}

	httpRequestsInflight := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "http_requests_inflight_total",
		Help:        "All the requests in progress",
		ConstLabels: labels,
	}, []string{"method", "path"})

	httpRequestsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        "http_requests_total",
		Help:        "Count of all HTTP requests",
		ConstLabels: labels,
	}, []string{"status_code", "method", "path"})

	httpRequestDuration := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        "http_request_duration_seconds",
		Help:        "Duration of all HTTP requests",
		ConstLabels: labels,
	}, []string{"status_code", "method", "path"})

	f.prom.MustRegister(httpRequestsInflight)
	f.prom.MustRegister(httpRequestsTotal)
	f.prom.MustRegister(httpRequestDuration)

	chRequestBegin := make(chan requestLabel, 100)
	chRequestFinish := make(chan requestLabel, 100)

	go func() {
		for {
			select {
			case label := <-chRequestBegin:
				httpRequestsInflight.With(prometheus.Labels{
					"method": label.method,
					"path":   label.path,
				}).Inc()
			case label := <-chRequestFinish:
				httpRequestsInflight.With(prometheus.Labels{
					"method": label.method,
					"path":   label.path,
				}).Dec()
				httpRequestsTotal.With(prometheus.Labels{
					"status_code": strconv.Itoa(label.statusCode),
					"method":      label.method,
					"path":        label.path,
				}).Inc()
				httpRequestDuration.With(prometheus.Labels{
					"status_code": strconv.Itoa(label.statusCode),
					"method":      label.method,
					"path":        label.path,
				}).Observe(label.duration.Seconds())
			}
		}
	}()

	f.server.Use(func(ctx *fiber.Ctx) error {
		// Ignore metrics path
		if ctx.Path() == metricsPath {
			return ctx.Next()
		}

		start := time.Now()
		method := ctx.Method()
		path := strings.Replace(ctx.Path(), "/", "", 1) // for some reason, the path is prefixed with a / will cause prometheus to fail. XD

		chRequestBegin <- requestLabel{
			method: method,
			path:   path,
		}
		defer func() {
			chRequestFinish <- requestLabel{
				statusCode: ctx.Response().StatusCode(),
				method:     method,
				path:       path,
				duration:   time.Since(start),
			}
		}()

		return ctx.Next()
	})

	f.server.Get(metricsPath, adaptor.HTTPHandler(promhttp.HandlerFor(f.prom, promhttp.HandlerOpts{})))
}
