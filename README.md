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

- /etc/fluent-agent.logs

```
www0  /tmp/www0_access.log
www1  /tmp/www1_access.log
www2  /tmp/www2_access.log
```

The following metrics are exported:

```
# HELP fluent_agent_lite_cpu_time fluent-agent-lite cpu time
# TYPE fluent_agent_lite_cpu_time gauge
fluent_agent_lite_cpu_time{id="www0"} 0.07
fluent_agent_lite_cpu_time{id="www1"} 0.08
fluent_agent_lite_cpu_time{id="www2"} 0.09
# HELP fluent_agent_lite_down the fluent-agent-lite processes down
# TYPE fluent_agent_lite_down gauge
fluent_agent_lite_down 0
# HELP fluent_agent_lite_exporter_scrape_failures_total Number of errors while scraping fluent-agent-lite.
# TYPE fluent_agent_lite_exporter_scrape_failures_total counter
fluent_agent_lite_exporter_scrape_failures_total 2
# HELP fluent_agent_lite_resident_memory_usage fluent-agent-lite resident memory usage
# TYPE fluent_agent_lite_resident_memory_usage gauge
fluent_agent_lite_resident_memory_usage{id="www0"} 1.0359488e+07
fluent_agent_lite_resident_memory_usage{id="www1"} 1.0459488e+07
fluent_agent_lite_resident_memory_usage{id="www2"} 1.0559488e+07
# HELP fluent_agent_lite_up the fluent-agent-lite processes up
# TYPE fluent_agent_lite_up gauge
fluent_agent_lite_up 3
# HELP fluent_agent_lite_up_expect the fluent-agent-lite processes expect
# TYPE fluent_agent_lite_up_expect gauge
fluent_agent_lite_up_expect 3
# HELP fluent_agent_lite_virtual_memory_usage fluent-agent-lite virtual memory usage
# TYPE fluent_agent_lite_virtual_memory_usage gauge
fluent_agent_lite_virtual_memory_usage{id="www0"} 6.3791872e+07
fluent_agent_lite_virtual_memory_usage{id="www1"} 7.3791872e+07
fluent_agent_lite_virtual_memory_usage{id="www2"} 8.3791872e+07
```

# command line options

Name     | Description | Default | note
---------|-------------|----|----
web.listen-address | Address on which to expose metrics and web interface | 9269 |
web.telemetry-path | Path under which to expose metrics | /metrics |
log.level | Log level | info |

# How to build

```
$ make
```
