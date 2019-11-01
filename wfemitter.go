// Package ginwavefront is a Gin middleware to emit metrics to Wavefront.
package ginwavefront

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	wavefront "github.com/wavefronthq/wavefront-sdk-go/senders"
)

const (
	// ErrCreateSender in case any errors occur while creating the Wavefront Direct Sender
	ErrCreateSender = "error creating wavefront sender: %s"
)

// WavefrontConfig configures the direct ingestion sender to Wavefront.
type WavefrontConfig struct {
	// Wavefront URL of the form https://<INSTANCE>.wavefront.com.
	Server string
	// Wavefront API token with direct data ingestion permission.
	Token string
	// Max batch of data sent per flush interval.
	BatchSize int
	// Max batch of data sent per flush interval.
	MaxBufferSize int
	// Interval (in seconds) at which to flush data to Wavefront.
	FlushInterval int
	// Map of Key-Value pairs (strings) associated with each data point sent to Wavefront.
	PointTags map[string]string
	// Name of the app that emits metrics.
	Source string
	// Prefix added to all metrics
	MetricPrefix string
}

// WavefrontEmitter creates a new direct sender to Wavefront and returns a handlerfunc
func WavefrontEmitter(w *WavefrontConfig) (gin.HandlerFunc, error) {
	dc := &wavefront.DirectConfiguration{
		Server:               w.Server,
		Token:                w.Token,
		BatchSize:            w.BatchSize,
		MaxBufferSize:        w.MaxBufferSize,
		FlushIntervalSeconds: w.FlushInterval,
	}

	sender, err := wavefront.NewDirectSender(dc)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateSender, err.Error())
	}

	if w.PointTags == nil {
		w.PointTags = make(map[string]string)
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)
		statusCode := c.Writer.Status()
		bytesOut := c.Writer.Size()
		bytesIn := c.Request.ContentLength

		// Add tags
		w.PointTags["path"] = c.Request.URL.Path
		w.PointTags["clientIP"] = c.ClientIP()
		w.PointTags["method"] = c.Request.Method
		w.PointTags["userAgent"] = c.Request.UserAgent()

		// Send metrics
		// <metricName> <metricValue> [<timestamp>] source=<source> [pointTags]
		sender.SendMetric(strings.Join([]string{w.MetricPrefix, ".latency"}, ""), float64(latency.Milliseconds()), end.Unix(), w.Source, w.PointTags)
		sender.SendMetric(strings.Join([]string{w.MetricPrefix, ".bytes.in"}, ""), float64(bytesIn), end.Unix(), w.Source, w.PointTags)
		sender.SendMetric(strings.Join([]string{w.MetricPrefix, ".bytes.out"}, ""), float64(bytesOut), end.Unix(), w.Source, w.PointTags)
		switch {
		case statusCode > 199 && statusCode < 300:
			sender.SendDeltaCounter(strings.Join([]string{w.MetricPrefix, ".status.success"}, ""), 1, w.Source, w.PointTags)
		case statusCode > 299 && statusCode < 400:
			sender.SendDeltaCounter(strings.Join([]string{w.MetricPrefix, ".status.redirection"}, ""), 1, w.Source, w.PointTags)
		case statusCode > 399 && statusCode < 500:
			sender.SendDeltaCounter(strings.Join([]string{w.MetricPrefix, ".status.error.client"}, ""), 1, w.Source, w.PointTags)
		case statusCode > 499 && statusCode < 600:
			sender.SendDeltaCounter(strings.Join([]string{w.MetricPrefix, ".status.error.server"}, ""), 1, w.Source, w.PointTags)
		}
	}, nil
}
