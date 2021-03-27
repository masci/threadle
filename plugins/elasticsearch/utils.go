package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/masci/threadle/intake"
)

func addDocument(indexer esutil.BulkIndexer, d *document) error {
	jsonData, err := json.Marshal(d)
	if err != nil {
		return err
	}

	return indexer.Add(
		context.Background(),
		esutil.BulkIndexerItem{
			Action: "index",
			Body:   bytes.NewReader(jsonData),
		},
	)
}

// getV1MetricDocument converts a Datadog metric into an ECS compatible document
func getV1MetricDocument(m *intake.V1Metric) *document {
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

	return &d
}

func getHostMetadataDocument(hm *intake.HostMeta) *document {
	d := document{}
	d["@timestamp"] = time.Now().Format(time.RFC3339)
	d["host"] = host{
		Name:         hm.Meta.Hostname,
		Hostname:     hm.Meta.Hostname,
		ID:           hm.UUID,
		Architecture: hm.SystemStats.Machine,
		Mac:          hm.Network.Mac,
		IP:           hm.Network.IP,
	}
	// Convert tags to labels
	if len(hm.HostTags.System) > 0 {
		labels := labels{}
		for _, t := range hm.HostTags.System {
			toks := strings.Split(t, ":")
			if len(toks) < 2 {
				toks = append(toks, "")
			}
			labels[toks[0]] = toks[1]
		}
		d["labels"] = labels
	}
	return &d
}

func getProcDocuments(snap *intake.ProcessSnapshot) []*document {
	docs := make([]*document, len(snap.ProcessList))
	ts := time.Unix(int64(snap.Timestamp), 0).Format(time.RFC3339)
	for _, p := range snap.ProcessList {
		d := document{}
		d["@timestamp"] = ts
		d["process"] = p
		docs = append(docs, &d)
	}
	return docs
}

func setupIndex(es *elasticsearch.Client, indexName string) (created bool, err error) {
	var r *esapi.Response

	// check index exists
	r, err = es.Indices.Exists([]string{indexName})
	if err != nil {
		return
	}
	defer r.Body.Close()

	// if it doesn't, create the index and setup the mapping
	if r.StatusCode == 404 {
		r, err = es.Indices.Create(
			indexName,
			// map any number to float because we don't know in advance
			// what data the Datadog Agent will send. Probably we can
			// do better than this but I needed something quick.
			es.Indices.Create.WithBody(strings.NewReader(`{
				"mappings": {
					"dynamic_templates": [
					  {
						"integers": {
						  "match_mapping_type": "long",
						  "mapping": {
							"type": "double"
						  }
						}
					  }
					]
				}
			}`)),
		)
		if err != nil {
			return
		}
		if r.IsError() {
			err = fmt.Errorf("Cannot create index: %s", r)
			return
		}
		r.Body.Close()
		created = true
		return
	} else if r.IsError() {
		err = fmt.Errorf("response error querying the index: %s", r)
		return
	}

	// nothing to do
	return
}
