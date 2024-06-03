package httputil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// UnmarshalResponseBody ...
func UnmarshalResponseBody(resp *http.Response, result interface{}) error {
	body, err := RawResponseBody(resp)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, result); err != nil {
		return err
	}

	return nil
}

// RawResponseBody ...
func RawResponseBody(resp *http.Response) ([]byte, error) {
	r1, r2, err := DrainBody(resp.Body)
	resp.Body = r2
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(r1)
}

func BuildErrorWithResponse(err error, resp *http.Response) error {
	raw, _ := RawResponseBody(resp)
	return fmt.Errorf("%w, statusCode: %d, body: %s", err, resp.StatusCode, string(raw))
}
