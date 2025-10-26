package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func MakeJSONRequest[R, E any](method, url string, data io.Reader, headers ...map[string]string) (okResp R, errResp E, err error) {
	response, err := MakeRequest(method, url, data, headers...)

	if err != nil {
		return
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	if response.StatusCode != http.StatusOK {
		err = decoder.Decode(&errResp)
	} else {
		err = decoder.Decode(&okResp)
	}

	return
}

func MakeRequest(method, url string, data io.Reader, headers ...map[string]string) (*http.Response, error) {
	if req, err := http.NewRequest(method, url, data); err != nil {
		return nil, fmt.Errorf("error create request: %s", err)

	} else {
		req.Header.Set("Content-Type", "application/json")

		for _, head := range headers {
			for k, v := range head {
				req.Header.Set(k, v)
			}
		}

		client := http.Client{}

		if resp, err := client.Do(req); err != nil {
			return nil, fmt.Errorf("error make request: %s", err)

		} else {
			return resp, nil
		}
	}
}
