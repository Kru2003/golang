package middlewares

import (
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/config"
	"go.uber.org/zap"
)

type Middleware struct {
	config config.AppConfig
	logger *zap.Logger
}

func NewMiddleware(cfg config.AppConfig, logger *zap.Logger) Middleware {
	return Middleware{
		config: cfg,
		logger: logger,
	}
}
