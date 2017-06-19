package main

import (
	"flag"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/log"
	"github.com/prometheus/procfs"
)

var (
	// command line parameters
	listenAddress = flag.String("web.listen-address", "9270", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
)

const (
	// Can't use '-' for the metric name
	namespace = "fluent_agent_lite"
)

type TargetLogConfig struct {
	tag         string
	logFilePath string
}

type Exporter struct {
	mutex sync.RWMutex

	scrapeFailures          prometheus.Counter
	cpuTime                 *prometheus.GaugeVec
	virtualMemory           *prometheus.GaugeVec
	residentMemory          *prometheus.GaugeVec
	fluentAgentLiteUpExpect prometheus.Gauge
	fluentAgentLiteUp       prometheus.Gauge
	fluentAgentLiteDown     prometheus.Gauge
}

func NewExporter() *Exporter {
	return &Exporter{
		scrapeFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_scrape_failures_total",
			Help:      "Number of errors while scraping fluent-agent-lite.",
		}),
		cpuTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cpu_time",
				Help:      "fluent-agent-lite cpu time",
			},
			[]string{"id"},
		),
		virtualMemory: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "virtual_memory_usage",
				Help:      "fluent-agent-lite virtual memory usage",
			},
			[]string{"id"},
		),
		residentMemory: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "resident_memory_usage",
				Help:      "fluent-agent-lite resident memory usage",
			},
			[]string{"id"},
		),
		fluentAgentLiteUpExpect: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up_expect",
			Help:      "the fluent-agent-lite processes expect",
		}),
		fluentAgentLiteUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "the fluent-agent-lite processes up",
		}),
		fluentAgentLiteDown: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "down",
			Help:      "the fluent-agent-lite processes down",
		}),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.scrapeFailures.Describe(ch)
	e.cpuTime.Describe(ch)
	e.virtualMemory.Describe(ch)
	e.residentMemory.Describe(ch)
	e.fluentAgentLiteUpExpect.Describe(ch)
	e.fluentAgentLiteUp.Describe(ch)
	e.fluentAgentLiteDown.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// To protect metrics from concurrent collects.
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.collect(ch)
}

func (e *Exporter) collect(ch chan<- prometheus.Metric) {
	e.fluentAgentLiteUp.Set(0)
	e.fluentAgentLiteDown.Set(0)

	targetLogConfigs, err := e.getTargetLogConfigs()
	if err != nil {
		e.scrapeFailures.Inc()
		e.scrapeFailures.Collect(ch)
		return
	}

	log.Debugf("targetLogConfigs = %v", targetLogConfigs)

	processes := 0
	for _, targetLogConfig := range targetLogConfigs {
		log.Debugf("fluent-agent-lite targetLogConfig = %v", targetLogConfig)

		procStat, err := e.getProcStat(targetLogConfig)
		if err != nil {
			e.cpuTime.WithLabelValues(targetLogConfig.tag).Set(0)
			e.virtualMemory.WithLabelValues(targetLogConfig.tag).Set(float64(0))
			e.residentMemory.WithLabelValues(targetLogConfig.tag).Set(float64(0))

			e.scrapeFailures.Inc()
			continue
		}

		e.cpuTime.WithLabelValues(targetLogConfig.tag).Set(procStat.CPUTime())
		e.virtualMemory.WithLabelValues(targetLogConfig.tag).Set(float64(procStat.VirtualMemory()))
		e.residentMemory.WithLabelValues(targetLogConfig.tag).Set(float64(procStat.ResidentMemory()))

		processes++
	}

	upExpect := len(targetLogConfigs)

	e.fluentAgentLiteUpExpect.Set(float64(upExpect))
	e.fluentAgentLiteUp.Set(float64(processes))
	e.fluentAgentLiteDown.Set(float64(upExpect - processes))

	e.scrapeFailures.Collect(ch)
	e.cpuTime.Collect(ch)
	e.virtualMemory.Collect(ch)
	e.residentMemory.Collect(ch)
	e.fluentAgentLiteUpExpect.Collect(ch)
	e.fluentAgentLiteUp.Collect(ch)
	e.fluentAgentLiteDown.Collect(ch)
}

func (e *Exporter) getTargetLogConfigs() ([]TargetLogConfig, error) {
	// see also
	// https://github.com/tagomoris/fluent-agent-lite/blob/v1.0/package/fluent-agent-lite.init
	command := `
		source /etc/fluent-agent-lite.conf

		lines=$(echo "$LOGS" | grep -v '^#' | grep -v '^$' | wc -l | awk '{print $1;}')

		for (( i = 0; i < $lines; i++ )); do
			lineno=$((i + 1))
			line=$(echo "$LOGS" | grep -v '^#' | tail -n $lineno | head -n 1)
			tag=$(echo $line | awk '{print $1;}')
			logFilePath=$(echo $line | awk '{print $2;}')

			echo "${tag}:${logFilePath}"
		done
		`

	out, err := exec.Command("/bin/bash", "-c", command).Output()
	if err != nil {
		log.Error(err)
		return []TargetLogConfig{}, err
	}

	var targetLogConfigs []TargetLogConfig
	lines := strings.TrimSpace(string(out))

	for _, line := range strings.Split(lines, "\n") {
		tokens := strings.Split(line, ":")
		targetLogConfigs = append(targetLogConfigs,
			TargetLogConfig{
				tag:         tokens[0],
				logFilePath: tokens[1],
			})
	}

	return targetLogConfigs, nil
}

func (e *Exporter) getProcStat(targetLogConfig TargetLogConfig) (procfs.ProcStat, error) {
	procfsPath := procfs.DefaultMountPoint
	fs, err := procfs.NewFS(procfsPath)
	if err != nil {
		log.Error(err)
		return procfs.ProcStat{}, err
	}

	targetPid, err := e.resolveTargetPid(targetLogConfig)
	if err != nil {
		return procfs.ProcStat{}, err
	}

	log.Debugf("targetPid = %v", targetPid)

	proc, err := fs.NewProc(targetPid)
	if err != nil {
		log.Error(err)
		return procfs.ProcStat{}, err
	}

	procStat, err := proc.NewStat()
	if err != nil {
		log.Error(err)
		return procfs.ProcStat{}, err
	}

	return procStat, nil
}

func (e *Exporter) resolveTargetPid(targetLogConfig TargetLogConfig) (int, error) {
	pgrepArg := targetLogConfig.tag + " " + targetLogConfig.logFilePath

	log.Debugf("pgrep arg = [%s]", pgrepArg)

	out, err := exec.Command("pgrep", "-n", "-f", pgrepArg).Output()
	if err != nil {
		log.Error(err)
		return 0, err
	}

	targetPid, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return targetPid, nil
}

func main() {
	flag.Parse()

	exporter := NewExporter()
	prometheus.MustRegister(exporter)

	log.Infof("Starting Server: %s", *listenAddress)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>fluent-agent-lite Exporter</title></head>
		<body>
		<h1>fluent-agent-lite Exporter v` + version + `</h1>
		<p><a href="` + *metricsPath + `">Metrics</a></p>
		</body>
		</html>`))
	})

	log.Fatal(http.ListenAndServe(":"+*listenAddress, nil))
}
