# fluent-agent-lite_exporter

[![CircleCI](https://circleci.com/gh/matsumana/fluent-agent-lite_exporter/tree/master.svg?style=shield)](https://circleci.com/gh/matsumana/fluent-agent-lite_exporter/tree/master)

[fluent-agent-lite](https://github.com/tagomoris/fluent-agent-lite) exporter for [Prometheus](https://prometheus.io/)

# export metrics

- fluent_agent_lite_cpu_time
- fluent_agent_lite_resident_memory_usage
- fluent_agent_lite_virtual_memory_usage
- fluent_agent_lite_up_expect
- fluent_agent_lite_up
- fluent_agent_lite_down

e.g.

```
# HELP fluent_agent_lite_cpu_time fluent-agent-lite cpu time
# TYPE fluent_agent_lite_cpu_time gauge
fluent_agent_lite_cpu_time{id="www0"} 0
fluent_agent_lite_cpu_time{id="www1"} 0
fluent_agent_lite_cpu_time{id="www2"} 0.09
# HELP fluent_agent_lite_down the fluent-agent-lite processes down
# TYPE fluent_agent_lite_down gauge
fluent_agent_lite_down 2
# HELP fluent_agent_lite_exporter_scrape_failures_total Number of errors while scraping fluent-agent-lite.
# TYPE fluent_agent_lite_exporter_scrape_failures_total counter
fluent_agent_lite_exporter_scrape_failures_total 2
# HELP fluent_agent_lite_resident_memory_usage fluent-agent-lite resident memory usage
# TYPE fluent_agent_lite_resident_memory_usage gauge
fluent_agent_lite_resident_memory_usage{id="www0"} 0
fluent_agent_lite_resident_memory_usage{id="www1"} 0
fluent_agent_lite_resident_memory_usage{id="www2"} 1.0559488e+07
# HELP fluent_agent_lite_up the fluent-agent-lite processes up
# TYPE fluent_agent_lite_up gauge
fluent_agent_lite_up 1
# HELP fluent_agent_lite_up_expect the fluent-agent-lite processes expect
# TYPE fluent_agent_lite_up_expect gauge
fluent_agent_lite_up_expect 3
# HELP fluent_agent_lite_virtual_memory_usage fluent-agent-lite virtual memory usage
# TYPE fluent_agent_lite_virtual_memory_usage gauge
fluent_agent_lite_virtual_memory_usage{id="www0"} 0
fluent_agent_lite_virtual_memory_usage{id="www1"} 0
fluent_agent_lite_virtual_memory_usage{id="www2"} 8.3791872e+07
```

# command line options

Name     | Description | Default | note
---------|-------------|----|----
web.listen-address | Address on which to expose metrics and web interface | 9270 |
web.telemetry-path | Path under which to expose metrics | /metrics |
log.level | Log level | info |

# How to build

```
$ make
```
