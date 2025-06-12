package controllers

import (
	pMetrics "git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/pkg/prometheus"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type MetricsController struct {
	logger   *zap.Logger
	pMetrics *pMetrics.PrometheusMetrics
}

func InitMetricsController(logger *zap.Logger, pMetrics *pMetrics.PrometheusMetrics) (*MetricsController, error) {
	return &MetricsController{
		logger:   logger,
		pMetrics: pMetrics,
	}, nil
}

func (mc *MetricsController) Metrics(ctx *fiber.Ctx) error {
	return adaptor.HTTPHandler(promhttp.Handler())(ctx)
}
