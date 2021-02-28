package intake

import "encoding/json"

// Point is an alias for an array of floats
type Point []float64

// V1Metric represents a datapoint series
type V1Metric struct {
	Metric         string
	Points         []Point
	Tags           []string
	Host           string
	Device         string
	Type           string
	Interval       int
	SourceTypeName string `json:"source_type_name"`
}

type series struct {
	Series []V1Metric `json:"series"`
}

// GetV1Metrics decodes a payload and returns a slice of Metric
func GetV1Metrics(payload []byte) ([]V1Metric, error) {
	s := series{[]V1Metric{}}

	return s.Series, json.Unmarshal(payload, &s)
}
