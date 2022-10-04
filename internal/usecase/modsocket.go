package usecase

import (
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
	"sync"
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

	successfullInternalModsocketConnections = promauto.NewCounter(prometheus.CounterOpts{
		Name: "successfull_internal_modsocket_connections",
		Help: "The total number of successfull internal modsocket connections",
	})

	failedInternalModsocketConnections = promauto.NewCounter(prometheus.CounterOpts{
		Name: "failed_internal_modsocket_connections",
		Help: "The total number of failed internal modsocket connections",
	})
)

func startMonitoring() {
	for {

		monitorModSockets()

		// wait 10 seconds
		time.Sleep(10 * time.Second)
	}
}

func monitorModSockets() {
	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		monitorInternalModsocket()
	}()
	go func() {
		defer wg.Done()
		monitorExternalModsocket()
	}()

	wg.Wait()
}

func monitorInternalModsocket() {

	working := checkModSocketConnection("wss", ModsocketUrl)

	if working {
		successfullInternalModsocketConnections.Inc()
	} else {
		log.Warn().Msg("modsocket internal connection failed")
		failedInternalModsocketConnections.Inc()
	}
}

func monitorExternalModsocket() {

	working := checkModSocketConnection("ws", internalModSocket())

	if working {
		successfullModsocketConnections.Inc()
	} else {
		log.Warn().Msg("modsocket external connection failed")
		failedModsocketConnections.Inc()
	}
}

func checkModSocketConnection(scheme, host string) bool {
	u := url.URL{Scheme: scheme, Host: host, Path: "/modsocket"}

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

func internalModSocket() string {
	v := os.Getenv("MODSOCKET_INTERNAL_URL")
	if v == "" {
		log.Panic().Msgf("MODSOCKET_INTERNAL_URL is not set")
	}
	return v
}
