package plugins

import (
	"regexp"
	"testing"

	"github.com/masci/threadle/intake"
	"github.com/stretchr/testify/require"
)

func TestExcludeV1Metrics(t *testing.T) {
	testcases := []struct {
		exclude  Filters
		metrics  []intake.V1Metric
		expected []intake.V1Metric
	}{
		{
			exclude: Filters{regexp.MustCompile(`.+\.datadog\..+`)},
			metrics: []intake.V1Metric{
				{
					Metric: "system.cpu.system",
				},
				{
					Metric: "foo.datadog.bar",
				},
			},
			expected: []intake.V1Metric{
				{
					Metric: "system.cpu.system",
				},
			},
		},
		{
			exclude: Filters{regexp.MustCompile(`.+`)},
			metrics: []intake.V1Metric{
				{
					Metric: "system.cpu.system",
				},
				{
					Metric: "foo.datadog.bar",
				},
			},
			expected: []intake.V1Metric{},
		},
		{
			exclude: Filters{regexp.MustCompile(`.+\.datadog\..+`)},
			metrics: []intake.V1Metric{
				{
					Metric: "system.cpu.system",
				},
				{
					Metric: "system.cpu.user",
				},
			},
			expected: []intake.V1Metric{
				{
					Metric: "system.cpu.system",
				},
				{
					Metric: "system.cpu.user",
				},
			},
		},
	}

	for _, testcase := range testcases {
		res := ExcludeV1Metrics(testcase.metrics, testcase.exclude)
		require.Equal(t, testcase.expected, res)
	}
}
