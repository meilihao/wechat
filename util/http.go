package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var DefaultHttpClient = http.DefaultClient

func PostJSON(baseURL string, reqestBody interface{}, response interface{}) (err error) {
	buf := bytes.NewBuffer(nil)

	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	if err = encoder.Encode(reqestBody); err != nil {
		return
	}

	resp, err := DefaultHttpClient.Post(baseURL, "application/json; charset=utf-8", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(response)
}
