package main

import (
	"log"
	"os"
	"net/url"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr     string
	endpoint *url.URL
)

func main() {
	loadEnv()

	c := newCollector(endpoint)
	prometheus.Register(c)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}

func loadEnv() {
	addr = getEnvOr("YARN_PROMETHEUS_LISTEN_ADDR", ":9113")

	scheme := getEnvOr("YARN_PROMETHEUS_ENDPOINT_SCHEME", "http")

	hostName, err := os.Hostname()
	log.Println(hostName)
	host := getEnvOr("YARN_PROMETHEUS_ENDPOINT_HOST", hostName)
	port := getEnvOr("YARN_PROMETHEUS_ENDPOINT_PORT", "23188")
	path := getEnvOr("YARN_PROMETHEUS_ENDPOINT_PATH", "ws/v1/cluster/metrics")

	e, err := url.Parse(scheme + "://" + host + ":" + port + "/" + path)
	if err != nil {
		log.Fatal()
	}

	endpoint = e
}

func getEnvOr(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}
