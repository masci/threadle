package intake

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetProcessSnapshots(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/host_meta.json")
	if err != nil {
		t.Fatalf("Error loading golden file: %s", err)
	}
	hostMeta, err := DecodeHostMeta(content)

	snapshots := hostMeta.GetProcessSnapshots()
	require.Len(t, snapshots, 1)

	processes := snapshots[0].ProcessList
	require.Len(t, processes, 20)
	require.Equal(t, "Google", processes[0].Name)
}
