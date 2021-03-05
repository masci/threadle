package plugins

import (
	"regexp"

	"github.com/masci/threadle/intake"
)

// Filters is a slice of exclusion filters
type Filters []*regexp.Regexp

// ExcludeV1Metrics drops metrics according to one or more exclusion filters for their name
func ExcludeV1Metrics(metrics []intake.V1Metric, exclude Filters) []intake.V1Metric {
	// filter in-place
	n := 0
	for _, m := range metrics {
		keep := true
		for _, reg := range exclude {
			if reg.MatchString(m.Metric) {
				keep = false
				break
			}
		}
		if keep {
			metrics[n] = m
			n++
		}
	}
	return metrics[:n]
}
