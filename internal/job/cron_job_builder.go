package job

import (
	"github.com/jasonzhao47/cuddle/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"strconv"
	"time"
)

type CronJobBuilder struct {
	l   logger.Logger
	vec *prometheus.SummaryVec
}

func NewCronJobBuilder(l logger.Logger, opts prometheus.SummaryOpts) *CronJobBuilder {
	vector := prometheus.NewSummaryVec(opts, []string{"job", "success"})
	return &CronJobBuilder{l: l, vec: vector}
}

func (c *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()
	return cronJobAdapter(func() {
		start := time.Now()
		err := job.Run()
		if err != nil {
			c.l.Error("error running job",
				logger.Error(err),
				logger.String("name", name))
		}
		c.l.Debug("job completed", logger.String("name", name))
		duration := time.Since(start)
		c.vec.WithLabelValues(name, strconv.FormatBool(err == nil)).Observe(float64(duration))
	})
}

type cronJobAdapter func()

func (c cronJobAdapter) Run() {
	c()
}
