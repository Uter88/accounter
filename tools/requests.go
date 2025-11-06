package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strings"
)

type Params struct {
	values url.Values
}

func NewParams() Params {
	return Params{
		values: make(url.Values),
	}
}

func (p *Params) Add(key string, value any) *Params {
	p.values.Add(key, fmt.Sprintf("%v", value))
	return p
}

func (p *Params) Set(key string, value any) *Params {
	p.values.Set(key, fmt.Sprintf("%v", value))
	return p
}

func (p *Params) Encode() string {
	if len(p.values) == 0 {
		return ""
	}

	return "?" + p.values.Encode()
}

type Request[R any] struct {
	method  string
	api     string
	path    string
	params  Params
	data    io.Reader
	headers map[string]string
}

func (p Request[R]) Method(method string) Request[R] {
	p.method = method

	return p
}

func (p Request[R]) Path(path string) Request[R] {
	if strings.HasPrefix(path, "/") {
		path = strings.Replace(path, "/", "", 1)
	}

	p.path = path

	return p
}

func (p Request[R]) Param(key string, value any) Request[R] {
	p.params.Set(key, value)

	return p
}

func (p Request[R]) Params(params Params) Request[R] {
	for k, v := range params.values {
		p.params.Add(k, v)
	}

	return p
}

func (p Request[R]) Data(data io.Reader) Request[R] {
	p.data = data

	return p
}

func (p Request[R]) Headers(headers map[string]string) Request[R] {
	maps.Copy(p.headers, headers)

	return p
}

func NewRequest[R any](api string) Request[R] {
	return Request[R]{
		api:     api,
		method:  http.MethodGet,
		params:  NewParams(),
		headers: make(map[string]string),
	}
}

func (p Request[R]) Do() (R, R, error) {
	api := fmt.Sprintf("%s/%s%s", p.api, p.path, p.params.Encode())

	return MakeJSONRequest[R, R](p.method, api, p.data, p.headers)
}

func MakeJSONRequest[R, E any](method, url string, data io.Reader, headers ...map[string]string) (okResp R, errResp E, err error) {
	response, err := MakeRequest(method, url, data, headers...)

	if err != nil {
		return
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	if response.StatusCode != http.StatusOK {
		decoder.Decode(&errResp)

		err = errors.New(response.Status)
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
