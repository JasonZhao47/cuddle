package ioc

import (
	"github.com/jasonzhao47/cuddle/internal/service/sms"
	"github.com/jasonzhao47/cuddle/internal/service/sms/prometheus"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
)

func InitPrometheusService(sms sms.Service) *prometheus.Decorator {
	opts := prometheus2.SummaryOpts{
		Namespace: "jason_zhao",
		Subsystem: "sms",
		Name:      "sms_service",
	}
	return prometheus.NewDecorator(sms, opts)

}
