[![test](https://github.com/masci/threadle/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/masci/threadle/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/masci/threadle)](https://goreportcard.com/report/github.com/masci/threadle)

# Threadle

Threadle can ingest metrics from a Datadog Agent and send them to a custom storage using different plugins.

A use case example would be using Elasticsearch to store the timeseries and having Grafana visualize data.

![Datadog](img/datadog.png)
![Grafana](img/grafana.png)

## Quickstart

Install threadle:

```bash
$ go get -u github.com/masci/threadle
go: downloading github.com/masci/threadle
```

By default Threadle discards all the messages coming from the Datadog Agent so you have to enable at
least one output plugin. Create a basic configuration file named `threadle.yaml`:

```yaml
plugins:
  logger:
```

Launch Threadle from the same directory containing the config file:

```bash
$ $GOPATH/bin/threadle -c
Initializing plugin: logger
```

By default Threadle listens on port `3060`, point the Datadog Agent there by adding the following to your
`datadog.yaml` configuration file:

```yaml
additional_endpoints:
  "http://localhost:3060": ""
```

Restart the Datadog Agent, you're all set.

## Plugins

Threadle is a small tool I built for myself so it doesn't offer much out of the box, but adding a plugin
shouldn't be hard.

### Logger

The logger plugin just prints the payload received from the Datadog Agent on `stderr` in JSON format. It
is mostly intended for debugging but it has an option to make the log lines [ECS](https://www.elastic.co/guide/en/ecs/current/index.html)
compatible, in case you want to send them straight to an index in Elasticsearch without additional setup.
The plugin only has one configuration option:

```yaml
logger:
  ecs_compatible: true
```

### Elasticsearch
