package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"

	"github.com/prometheus/log"
)

// e2e test
func TestE2EPartialFailure(t *testing.T) {
	metrics, err := get("http://localhost:9269/metrics")
	if err != nil {
		t.Error("HttpClient.Get = %v", err)
	}

	log.Info(metrics)

	// fluent_agent_lite_cpu_time
	if !regexp.MustCompile(`fluent_agent_lite_cpu_time\{id="www0"\} 0`).MatchString(metrics) {
		t.Error(`fluent_agent_lite_cpu_time{id="www0"} doesn't match`)
	}
	if !regexp.MustCompile(`fluent_agent_lite_cpu_time\{id="www1"\} 0`).MatchString(metrics) {
		t.Error(`fluent_agent_lite_cpu_time{id="www1"} doesn't match`)
	}
	if !regexp.MustCompile(`fluent_agent_lite_cpu_time\{id="www2"\} 0`).MatchString(metrics) {
		t.Error(`fluent_agent_lite_cpu_time{id="www2"} doesn't match`)
	}

	// fluent_agent_lite_resident_memory_usage
	if !regexp.MustCompile(`fluent_agent_lite_resident_memory_usage\{id="www0"\} `).MatchString(metrics) {
		t.Error(`fluent_agent_lite_resident_memory_usage{id="www0"} doesn't match`)
	}
	if !regexp.MustCompile(`fluent_agent_lite_resident_memory_usage\{id="www1"\} `).MatchString(metrics) {
		t.Error(`fluent_agent_lite_resident_memory_usage{id="www1"} doesn't match`)
	}
	if !regexp.MustCompile(`fluent_agent_lite_resident_memory_usage\{id="www2"\} `).MatchString(metrics) {
		t.Error(`fluent_agent_lite_resident_memory_usage{id="www2"} doesn't match`)
	}

	// fluent_agent_lite_virtual_memory_usage
	if !regexp.MustCompile(`fluent_agent_lite_virtual_memory_usage\{id="www0"\} `).MatchString(metrics) {
		t.Error(`fluent_agent_lite_virtual_memory_usage{id="www0"} doesn't match`)
	}
	if !regexp.MustCompile(`fluent_agent_lite_virtual_memory_usage\{id="www1"\} `).MatchString(metrics) {
		t.Error(`fluent_agent_lite_virtual_memory_usage{id="www1"} doesn't match`)
	}
	if !regexp.MustCompile(`fluent_agent_lite_virtual_memory_usage\{id="www2"\} `).MatchString(metrics) {
		t.Error(`fluent_agent_lite_virtual_memory_usage{id="www2"} doesn't match`)
	}

	// fluent_agent_lite_exporter_scrape_failures_total
	if !regexp.MustCompile(`fluent_agent_lite_exporter_scrape_failures_total `).MatchString(metrics) {
		t.Error(`fluent_agent_lite_exporter_scrape_failures_total doesn't match`)
	}

	// fluent_agent_lite_up_expect
	if !regexp.MustCompile("fluent_agent_lite_up_expect 3").MatchString(metrics) {
		t.Error("fluent_agent_lite_up_expect doesn't match")
	}

	// fluent_agent_lite_up
	if !regexp.MustCompile("fluent_agent_lite_up 1").MatchString(metrics) {
		t.Error("fluent_agent_lite_up doesn't match")
	}

	// fluent_agent_lite_down
	if !regexp.MustCompile("fluent_agent_lite_down 2").MatchString(metrics) {
		t.Error("fluent_agent_lite_down doesn't match")
	}
}

func get(url string) (string, error) {
	log.Info(url)

	response, err := http.Get(url)
	if err != nil {
		log.Error("http.Get = %v", err)
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("ioutil.ReadAll = %v", err)
		return "", err
	}
	if response.StatusCode != 200 {
		log.Error("response.StatusCode = %v", response.StatusCode)
		return "", err
	}

	metrics := string(body)

	return metrics, nil
}
