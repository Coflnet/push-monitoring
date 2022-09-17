package usecase

import (
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"net/url"
	"time"
)

const (
	ModsocketUrl = "sky.coflnet.com"
)

var (
	successfullModsocketConnections = promauto.NewCounter(prometheus.CounterOpts{
		Name: "successfull_modsocket_connections",
		Help: "The total number of successfull modsocket connections",
	})

	failedModsocketConnections = promauto.NewCounter(prometheus.CounterOpts{
		Name: "failed_modsocket_connections",
		Help: "The total number of failed modsocket connections",
	})
)

func monitorModSocket() {
	for {
		working := checkModSocketConnection()

		if working {
			log.Debug().Msg("modsocket is working")
			successfullModsocketConnections.Inc()
		} else {
			log.Warn().Msg("modsocket connection failed")
			failedModsocketConnections.Inc()
		}

		// wait 10 seconds
		time.Sleep(10 * time.Second)
	}
}

func checkModSocketConnection() bool {
	u := url.URL{Scheme: "wss", Host: ModsocketUrl, Path: "/modsocket"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to modsocket")
		return false
	}

	// close connection again
	err = c.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to close connection")
		return false
	}

	return true
}
