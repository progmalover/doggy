package doggy

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/negroni"
)

var (
	requestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_count",
		Help: "http request count.",
	}, []string{"code", "path"})

	requestLatencyHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "http_request_latency_histogram",
		Help: "http request latency histogram.",
	})
)

func Prometheus(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	now := time.Now()

	next(rw, r)

	status := 0
	if ww, ok := rw.(negroni.ResponseWriter); ok {
		status = ww.Status()
	}
	elasped := time.Now().Sub(now).Seconds()
	requestCounter.WithLabelValues(strconv.Itoa(status), r.URL.Path).Inc()
	requestLatencyHistogram.Observe(elasped)
}

func init() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestLatencyHistogram)
}
