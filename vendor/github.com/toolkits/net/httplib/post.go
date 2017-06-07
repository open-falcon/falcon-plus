package httplib

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func PostJSON(url string, v interface{}) (response []byte, err error) {
	var bs []byte
	bs, err = json.Marshal(v)
	if err != nil {
		return
	}

	bf := bytes.NewBuffer(bs)

	var resp *http.Response
	resp, err = http.Post(url, "application/json", bf)
	if err != nil {
		return
	}

	if resp.Body != nil {
		defer resp.Body.Close()
		response, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
	}

	if resp.StatusCode != 200 {
		err = errors.New("status code not equals 200")
	}

	return
}
