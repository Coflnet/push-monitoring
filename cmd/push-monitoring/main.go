package main

import (
	"github.com/Coflnet/push-monitoring/internal/metrics"
	"github.com/Coflnet/push-monitoring/internal/usecase"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msgf("starting app..")

	usecase.StartMonitoringServices()

	// start server
	err := metrics.StartServer()
	log.Panic().Err(err).Msg("failed to start server")
}
