package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	prometheusVersion "github.com/prometheus/common/version"
	"log"
	"net/http"
	"os"

	"random_exporter/collector"
)

const (
	namespace = "random"
)

var (
	int10 = prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace, "int", "int10"),
		"a random number from 0 to 10",
		nil,
		nil,
	)
)

type Exporter struct {
	Client *collector.RandomClient
}

func NewExporter() (*Exporter, error) {
	log.Println(int10)

	exp := Exporter{
		Client: collector.NewClient(42),
	}

	return &exp, nil
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- int10
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	if e.Client == nil {
		log.Println("shit happend")
	}

	metrics := e.Client.GetMetrics()
	ch <- prometheus.MustNewConstMetric(int10, prometheus.GaugeValue, metrics["int10"])
}

func init() {
	prometheus.MustRegister(prometheusVersion.NewCollector("random_exporter"))
}

func main() {
	var (
		showVersion   = flag.Bool("version", false, "print version information")
		listenAddress = flag.String("web.listen.address", ":9192", "the address it listens to")
		metricsPath   = flag.String("web.telemetry.path", "/metrics", "path under which metrics are exported")
	)
	flag.Parse()

	if *showVersion {
		fmt.Println("random value exporter version " + "0.0.1")
	}

	exporter, err := NewExporter()
	if err != nil {
		log.Fatal("no exporter")
		os.Exit(1)
	}

	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
					<head><title>random exporter for prometheus</title></head>
						<body>
							<a href='` + *metricsPath + `'>the metrics are here.</a>
						</body>
					</html>`))
	})

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
