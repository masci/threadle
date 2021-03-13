package intake

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeHostMeta(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/host_meta.json")
	if err != nil {
		t.Fatalf("Error loading golden file: %s", err)
	}

	_, err = DecodeHostMeta(content)
	// We don't do any logic there, let's just test unmarshalling works
	require.Nil(t, err)
}
