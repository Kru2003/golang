package cli

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/config"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"

	pMetrics "git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/pkg/prometheus"
)

// GetAPICommandDef runs app
func GetAPICommandDef(cfg config.AppConfig, logger *zap.Logger) cobra.Command {
	apiCommand := cobra.Command{
		Use:   "api",
		Short: "To start api",
		Long:  `To start api`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// Create fiber app
			app := fiber.New(fiber.Config{})

			promMetrics := pMetrics.InitPrometheusMetrics()

			// setup routes
			err := routes.Setup(app, logger, cfg, promMetrics)

			if err != nil {
				return err
			}

			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				if err := app.Listen(cfg.Port); err != nil {
					logger.Panic(err.Error())
				}
			}()

			<-interrupt
			logger.Info("gracefully shutting down...")
			if err := app.Shutdown(); err != nil {
				logger.Panic("error while shutdown server", zap.Error(err))
			}

			logger.Info("server stopped to receive new requests or connection.")
			return nil
		},
	}

	return apiCommand
}
