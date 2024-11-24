package monitoring

import (
	"expvar"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	requestCount   = expvar.NewInt("request_count")
	errorCount     = expvar.NewInt("error_count")
	totalLatency   = expvar.NewFloat("total_latency")
	averageLatency = expvar.NewFloat("average_latency")
)

var latencySum int64

func CollectExpvarMetrics(router *gin.Engine) {
	router.GET("/debug/vars", gin.WrapH(http.DefaultServeMux))
}

func UpdateMetrics(duration time.Duration, status int) {
	requestCount.Add(1)
	atomic.AddInt64(&latencySum, int64(duration.Milliseconds()))
	totalLatency.Set(float64(latencySum) / 1000)

	if requestCount.Value() > 0 {
		averageLatency.Set(float64(latencySum) / float64(requestCount.Value()))
	}

	if status >= 400 {
		errorCount.Add(1)
	}
}
