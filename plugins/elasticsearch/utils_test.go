package elasticsearch

import (
	"io/ioutil"
	"testing"

	"github.com/masci/threadle/intake"
	"github.com/stretchr/testify/require"
)

func TestGetV1MetricDocument(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/metrics.json")
	if err != nil {
		t.Fatalf("Error loading golden file: %s", err)
	}
	metrics, _ := intake.DecodeV1Metrics(content)

	expected := &document{
		"@timestamp": "2021-02-09T21:35:02Z",
		"host": host{
			Name:     "MacLastic2.local",
			Hostname: "MacLastic2.local",
		},
		"source_type_name":  "System",
		"system.cpu.system": 3.0419201600117516,
		"type":              "gauge",
	}

	doc := getV1MetricDocument(&(metrics[0]))
	require.Equal(t, expected, doc)
}
