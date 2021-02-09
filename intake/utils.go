package intake

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
	"net/http"
)

func inflate(in []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}

func readRequestBody(r *http.Request) (body []byte, err error) {
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}

	return inflate(body)
}
