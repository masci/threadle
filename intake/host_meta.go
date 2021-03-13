package intake

import "encoding/json"

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
		IP   string `json:"ipaddress"`
		IPV6 string `json:"ipaddressv6"`
		Mac  string `json:"macaddress"`
	}
	Logs      struct{ Transport string }
	Resources struct {
		Processes struct {
			Snaps []interface{}
		}
	}
}

// DecodeHostMeta decodes a HostMeta payload
func DecodeHostMeta(payload []byte) (*HostMeta, error) {
	hm := HostMeta{}
	err := json.Unmarshal(payload, &hm)
	return &hm, err
}
