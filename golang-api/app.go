// Golang API.
//
//	Schemes: http
//	Version: 0.0.1-alpha
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"time"

	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/cli"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/config"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/logger"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/routinewrapper"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func main() {
	// Collecting config from env or file or flag
	cfg := config.GetConfig()

	logger, err := logger.NewRootLogger(cfg.Debug, cfg.IsDevelopment)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	// this function will logged error log in sentry
	sentryLoggedFunc := func() {
		err := recover()

		if err != nil {
			sentry.CurrentHub().Recover(err)
			sentry.Flush(time.Second * 2)
		}
	}

	// routine wrapper will handle go routine error also an log into sentry
	routinewrapper.Init(sentryLoggedFunc)
	defer sentryLoggedFunc()

	err = cli.Init(cfg, logger)
	if err != nil {
		panic(err)
	}

}
