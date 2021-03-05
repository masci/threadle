package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/masci/threadle/intake"
	"github.com/masci/threadle/plugins"
	"github.com/spf13/viper"
)

var (
	es *elasticsearch.Client
)

// ES document
type document map[string]interface{}

// ECS compatible host field to be added to each document
type host struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
}

// ECS compatible labels field, we'll use it to store Datadog tags
type labels map[string]string

// Plugin implements plugins.Plugin
type Plugin struct{}

// Start subscribes to and processes the messages for the supported topics
func (*Plugin) Start(b *intake.PubSub) {
	// Create the Elasticsearch client
	var err error
	if es, err = elasticsearch.NewClient(elasticsearch.Config{
		// Auth
		CloudID:  viper.GetString("plugins.elasticsearch.cloudid"),
		Username: viper.GetString("plugins.elasticsearch.username"),
		Password: viper.GetString("plugins.elasticsearch.password"),
	}); err != nil {
		log.Fatalf("Error creating elasticsearch client: %s", err)
	}

	// Configure exclusion filters
	exclude := plugins.Filters{
		regexp.MustCompile(`.+\.datadog\..+`), // datadog.*
	}

	// Subcsribe to metrics messages
	go func() {
		var metrics []intake.V1Metric
		var err error
		for msg := range b.Subscribe(intake.SeriesEndpointV1) {
			metrics, err = intake.GetV1Metrics([]byte(msg))
			if err != nil {
				log.Println("error processing metrics: ", err)
				continue
			}

			process(metrics, exclude)
		}
	}()
}

// process reads all the metrics, build the corresponding ES documents and sends them
// using the _bulk api
func process(metrics []intake.V1Metric, exclude plugins.Filters) {
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  viper.GetString("plugins.elasticsearch.index"),
		Client: es,
	})
	if err != nil {
		log.Printf("Error creating the indexer: %s", err)
	}

	for _, m := range plugins.ExcludeV1Metrics(metrics, exclude) {
		d := document{}
		d["@timestamp"] = time.Unix(int64(m.Points[0][0]), 0).Format(time.RFC3339)
		d[m.Metric] = m.Points[0][1]
		if len(m.Tags) > 0 {
			labels := labels{}
			for _, t := range m.Tags {
				toks := strings.Split(t, ":")
				if len(toks) < 2 {
					toks = append(toks, "")
				}
				labels[toks[0]] = toks[1]
			}
			d["labels"] = labels
		}
		d["host"] = host{
			Name:     m.Host,
			Hostname: m.Host,
		}
		if m.Interval > 0 {
			d["interval"] = m.Interval
		}
		if m.Device != "" {
			d["device"] = m.Device
		}
		d["type"] = m.Type
		if m.SourceTypeName != "" {
			d["source_type_name"] = m.SourceTypeName
		}

		var jsonData []byte
		var err error
		if jsonData, err = json.Marshal(d); err != nil {
			log.Println(err)
			continue
		}

		err = indexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action: "index",
				Body:   bytes.NewReader(jsonData),
			},
		)
		if err != nil {
			log.Printf("Unexpected error: %s", err)
		}
	}

	if err := indexer.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	log.Println("flushed", indexer.Stats().NumFlushed, "created", indexer.Stats().NumCreated, "failed", indexer.Stats().NumFailed)
}
