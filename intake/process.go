package intake

// ProcessSnapshot represents a list of snapshots
type ProcessSnapshot struct {
	Timestamp   float64
	ProcessList []*Process
}

// Process represent a process stats
type Process struct {
	Username string
	CPUPct   float64
	MemPct   float64
	VMS      float64
	RSS      float64
	Name     string
	Pids     float64
}

// NewProcessSnapshot creates a ProcessSnapshot from raw data.
//
// The payload is a slice in the form [timestamp, process_list]
// where process_list is in turn an array of 5 items.
func NewProcessSnapshot(rawSnapshot interface{}) *ProcessSnapshot {
	processList := []*Process{}

	// interface{} --> [timestamp, process_list]
	snapshot := rawSnapshot.([]interface{})
	// process_list --> [process1, process2, ...]
	plist := snapshot[1].([]interface{})

	for _, p := range plist {
		// process --> [username, cpu_pct, ...]
		processList = append(processList, newProcess(p.([]interface{})))
	}
	return &ProcessSnapshot{
		Timestamp:   snapshot[0].(float64),
		ProcessList: processList,
	}
}

func newProcess(data []interface{}) *Process {
	return &Process{
		data[0].(string),
		data[1].(float64),
		data[2].(float64),
		data[3].(float64),
		data[4].(float64),
		data[5].(string),
		data[6].(float64),
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
