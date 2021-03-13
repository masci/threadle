package intake

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetV1Metrics(t *testing.T) {
	testcases := []struct {
		name    string
		metrics []V1Metric
		golden  string
	}{
		{
			name: "One Metric",
			metrics: []V1Metric{
				{
					Metric: "system.cpu.system",
					Points: []Point{
						{1612906502, 3.0419201600117516},
					},
					Tags:           []string{},
					Host:           "MacLastic2.local",
					Type:           "gauge",
					Interval:       0,
					SourceTypeName: "System",
				},
			},
			golden: "metrics",
		},
		{
			name:    "No Metrics",
			metrics: []V1Metric{},
			golden:  "metrics_empty",
		},
	}

	for _, testcase := range testcases {
		// get the raw payload
		content, err := ioutil.ReadFile("testdata/" + testcase.golden + ".json")
		if err != nil {
			t.Fatalf("Error loading golden file: %s", err)
		}
		got, err := DecodeV1Metrics(content)
		require.Nil(t, err)
		require.Equal(t, len(testcase.metrics), len(got))
		require.Equal(t, testcase.metrics, got)
	}
}
