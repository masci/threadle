package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/masci/threadle/intake"
	"github.com/masci/threadle/output"
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
	Name         string
	Hostname     string
	ID           string
	Architecture string
	Mac          string
	IP           string
}

// ECS compatible labels field, we'll use it to store Datadog tags
type labels map[string]string

// Plugin implements plugins.Plugin
type Plugin struct{}

// Start subscribes to and processes the messages for the supported topics
func (*Plugin) Start(b *intake.PubSub) {
	// Create the Elasticsearch client
	var err error
	cfg := getEsConfig()
	if es, err = elasticsearch.NewClient(*cfg); err != nil {
		output.FATAL.Fatalf("Error creating elasticsearch client: %s", err)
	}

	// Create the index if needed
	indexName := viper.GetString("plugins.elasticsearch.index")
	created, err := setupIndex(es, indexName)
	if err != nil {
		output.FATAL.Fatalf("Unexpected error setting up the index '%s': %s", indexName, err)
	}
	if created {
		output.INFO.Println("Index created:", indexName)
	}

	output.INFO.Println("Sending data to index:", indexName)

	// Configure exclusion filters for metrics
	exclude := plugins.GetFilters(viper.GetStringSlice("plugins.elasticsearch.exclude_metrics"))

	// Subcsribe to metrics messages
	go func() {
		for msg := range b.Subscribe(intake.SeriesEndpointV1) {
			metrics, err := intake.DecodeV1Metrics([]byte(msg))
			if err != nil {
				output.ERROR.Println("error processing metrics: ", err)
				continue
			}
			processV1Metrics(metrics, exclude)
		}
	}()

	// Subscribe to host metadata messages
	go func() {
		for msg := range b.Subscribe(intake.IntakeEndpointV1) {
			hostMeta, err := intake.DecodeHostMeta([]byte(msg))
			if err != nil {
				output.ERROR.Println("error processing host metadata: ", err)
				continue
			}
			processHostMeta(hostMeta)
		}
	}()
}

// processV1Metrics reads all the metrics, build the corresponding ES documents and stores them
// using the _bulk api
func processV1Metrics(metrics []intake.V1Metric, exclude plugins.Filters) {
	// Create the ES bulk indexer
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  viper.GetString("plugins.elasticsearch.index"),
		Client: es,
	})
	if err != nil {
		output.ERROR.Printf("Error creating the indexer: %s", err)
	}

	// Convert all the metrics and add them to the indexer
	for _, m := range plugins.ExcludeV1Metrics(metrics, exclude) {
		jsonData, err := json.Marshal(getV1MetricDocument(&m))
		if err != nil {
			output.ERROR.Println(err)
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
			output.ERROR.Printf("Unexpected error: %s", err)
		}
	}

	if err := indexer.Close(context.Background()); err != nil {
		output.FATAL.Fatalf("Unexpected error: %s", err)
	}

	output.DEBUG.Println("flushed", indexer.Stats().NumFlushed, "created", indexer.Stats().NumCreated, "failed", indexer.Stats().NumFailed)
}

func processHostMeta(hm *intake.HostMeta) {
	// Create the ES bulk indexer
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  viper.GetString("plugins.elasticsearch.index"),
		Client: es,
	})
	if err != nil {
		output.ERROR.Printf("Error creating the indexer: %s", err)
	}

	// Build the metadata document and add it to the bulk indexer
	if err = addDocument(indexer, getHostMetadataDocument(hm)); err != nil {
		output.ERROR.Printf("Error adding host meta to the indexer: %s", err)
	}

	// Go through all the snapshosts
	for _, snap := range hm.GetProcessSnapshots() {
		// Get the list of running processes
		for _, doc := range getProcDocuments(snap) {
			if err = addDocument(indexer, doc); err != nil {
				output.ERROR.Printf("Error adding snap to the indexer: %s", err)
				continue
			}
		}
	}

	// Flush data
	if err := indexer.Close(context.Background()); err != nil {
		output.FATAL.Fatalf("Unexpected error: %s", err)
	}

	output.DEBUG.Println("meta flushed", indexer.Stats().NumFlushed, "created", indexer.Stats().NumCreated, "failed", indexer.Stats().NumFailed)
}

func getEsConfig() *elasticsearch.Config {
	cfg := elasticsearch.Config{
		Username: viper.GetString("plugins.elasticsearch.username"),
		Password: viper.GetString("plugins.elasticsearch.password"),
	}

	// cloud id takes precedence over addresses
	if cloudid := viper.GetString("plugins.elasticsearch.cloudid"); cloudid != "" {
		cfg.CloudID = cloudid
		output.DEBUG.Println("Using Cloud ID", cloudid)
	} else {
		cfg.Addresses = viper.GetStringSlice("plugins.elasticsearch.addresses")
		output.DEBUG.Println("Using ES addresses", cfg.Addresses)
	}

	return &cfg
}
