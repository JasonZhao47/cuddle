package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

// 按每个请求接口分类
// 分类查看是对了还是错了

type Builder struct {
	Namespace  string
	Subsystem  string
	Name       string
	InstanceId string
}

func (b *Builder) BuildResponseTime() gin.HandlerFunc {
	labels := []string{"method", "pattern", "status"}
	// vector对应的不同的观测指标
	vec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Subsystem: b.Subsystem,
			Namespace: b.Namespace,
			Name:      b.Name + "_resp_time",
			ConstLabels: map[string]string{
				"instance_id": b.InstanceId,
			},
			Objectives: map[float64]float64{
				0.5:   0.01,
				0.7:   0.01,
				0.9:   0.01,
				0.99:  0.001,
				0.999: 0.0001,
			},
		},
		labels,
	)
	prometheus.MustRegister(vec)
	return func(ctx *gin.Context) {
		// 每个请求上报Prometheus的逻辑
		// 记录请求开始和请求结束的时间
		start := time.Now()
		defer func() {
			duration := time.Since(start).Milliseconds()
			method := ctx.Request.Method
			pattern := ctx.FullPath()
			status := ctx.Writer.Status()
			// Observe了啥纵坐标就是啥
			vec.WithLabelValues(method, pattern, strconv.Itoa(status)).
				Observe(float64(duration))
		}()
		ctx.Next()
	}
}

func (b *Builder) BuildActiveRequests() gin.HandlerFunc {
	// vector对应的不同的观测指标
	gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Subsystem: b.Subsystem,
			Namespace: b.Namespace,
			Name:      b.Name + "_active_reqs",
			ConstLabels: map[string]string{
				"instance_id": b.InstanceId,
			},
		},
	)
	prometheus.MustRegister(gauge)
	return func(ctx *gin.Context) {
		// 每个请求上报Prometheus的逻辑
		// 记录请求开始和请求结束的时间
		gauge.Inc()
		defer gauge.Desc()
		ctx.Next()
	}
}
