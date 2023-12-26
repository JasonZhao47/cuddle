package prometheus

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/service/sms"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

// 无侵入的Decorator

type Decorator struct {
	sms    sms.Service
	vector *prometheus.SummaryVec
}

func NewDecorator(sms sms.Service, opts prometheus.SummaryOpts) *Decorator {
	return &Decorator{
		sms: sms,
		vector: prometheus.NewSummaryVec(
			opts,
			[]string{"tpl_id"},
		),
	}
}

func (d *Decorator) Send(ctx context.Context, tplId string, args []string, phoneNums []string) error {
	start := time.Now()
	// 利用defer语义，在send之后执行持续时间的观测
	defer func() {
		duration := time.Since(start).Milliseconds()
		d.vector.WithLabelValues(tplId).Observe(float64(duration))
	}()
	return d.sms.Send(ctx, tplId, args, phoneNums)
}
