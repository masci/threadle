package intake

import (
	"encoding/json"
)

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

// HostMeta represents a metadata payload
type HostMeta struct {
	APIKey           string
	AgentVersion     string
	UUID             string
	InternalHostname string
	Os               string
	AgentFlavor      string `json:"agent-flavor"`
	Python           string
	SystemStats      struct {
		CPUCores  int
		Machine   string
		Platform  string
		PythonV   string
		Processor string
		MacV      interface{}
		NixV      interface{}
		FsdV      interface{}
		WinV      interface{}
	}
	Meta struct {
		SocketHostname string `json:"socket-hostname"`
		Timezones      []string
		SocketFqdn     string `json:"socket-fqdn"`
		Ec2Hostname    string `json:"ec2-hostname"`
		Hostname       string
		HostAliases    []string `json:"host_aliases"`
		InstanceID     string   `json:"instance-id"`
	}
	HostTags struct {
		System              []string `json:"system"`
		GoogleCloudPlatform []string `json:"google cloud platform,omitempty"`
	} `json:"host-tags"`
	Network struct {
		ID         string `json:"network-id"`
		PublicIPv4 string `json:"public-ipv4,omitempty"`
	}
	Logs      struct{ Transport string }
	Resources struct {
		Processes struct {
			Snaps []interface{}
		}
	}
}

// GetProcessSnapshots extracts the list of processes from the metadata payload
func (hm *HostMeta) GetProcessSnapshots() []*ProcessSnapshot {
	ret := []*ProcessSnapshot{}
	for _, rawSnap := range hm.Resources.Processes.Snaps {
		s := NewProcessSnapshot(rawSnap)
		ret = append(ret, s)
	}

	return ret
}

// DecodeV1Metrics decodes a payload and returns a slice of Metric
func DecodeV1Metrics(payload []byte) ([]V1Metric, error) {
	s := series{[]V1Metric{}}

	return s.Series, json.Unmarshal(payload, &s)
}

// DecodeHostMeta decodes a HostMeta payload
func DecodeHostMeta(payload []byte) (*HostMeta, error) {
	hm := HostMeta{}
	err := json.Unmarshal(payload, &hm)
	return &hm, err
}
