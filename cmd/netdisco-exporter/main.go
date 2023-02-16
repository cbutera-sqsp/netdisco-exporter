package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/cbutera-sqsp/netdisco-exporter/internal/exporter"
)

const (
	_pollDuration = 15 * time.Second
)

var (
	netdiscoHost     = os.Getenv("NETDISCO_HOST")
	netdiscoUsername = os.Getenv("NETDISCO_USERNAME")
	netdiscoPassword = os.Getenv("NETDISCO_PASSWORD")

	up = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "up",
		Help: "Netdisco-exporter status",
	})
)

func main() {
	// Checking for necessary ENV vars
	if netdiscoHost == "" {
		log.Fatalln("Please provide netdisco host url with http:// or https:// prefix via env var NETDISCO_HOST")
	}
	if netdiscoUsername == "" {
		log.Fatalln("Please provide netdisco login username via env var NETDISCO_USERNAME")
	}
	if netdiscoPassword == "" {
		log.Fatalln("Please provide netdisco login password via env var NETDISCO_PASSWORD")
	}

	// Starting Prometheus endpoint
	log.Info("Starting netdisco-exporter for netdisco host: ", netdiscoHost)
	go func() {
		log.Info("Started Prometheus metric endpoint")
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	up.Set(1)

	// Collect metrics based on ticker poll duration
	t := time.NewTicker(_pollDuration)
	for range t.C {
		exporter.ApiStatus()
		exporter.SearchDevice()
	}

}
