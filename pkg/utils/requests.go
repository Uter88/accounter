package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strings"
)

// Params request params
type Params struct {
	values url.Values
}

// NewParams creates new Params
func NewParams() Params {
	return Params{
		values: make(url.Values),
	}
}

// Merge pair of Params
func (p *Params) Merge(op Params) *Params {
	for k, v := range op.values {
		for i := range v {
			p.values.Add(k, v[i])
		}
	}

	return p
}

// Add add value to multiple key
func (p *Params) Add(key string, value any) *Params {
	p.values.Add(key, fmt.Sprintf("%v", value))
	return p
}

// Set replace value by specified key
func (p *Params) Set(key string, value any) *Params {
	p.values.Set(key, fmt.Sprintf("%v", value))
	return p
}

// Get returns value or values joined by comma by specified key
func (p *Params) Get(key string) string {
	if vals := p.values[key]; len(vals) > 0 {
		return strings.Join(vals, ",")
	}

	return p.values.Get(key)
}

// Encode encode values to query string starts with ?
func (p *Params) Encode() string {
	if len(p.values) == 0 {
		return ""
	}

	return "?" + p.values.Encode()
}

// Request HTTP requester
type Request[R any] struct {
	method  string
	api     string
	path    string
	params  Params
	data    io.Reader
	headers map[string]string
}

// Method set request method
func (p Request[R]) Method(method string) Request[R] {
	p.method = method

	return p
}

// Path set request path
func (p Request[R]) Path(path string) Request[R] {
	if strings.HasPrefix(path, "/") {
		path = strings.Replace(path, "/", "", 1)
	}

	p.path = path

	return p
}

// Param set request param
func (p Request[R]) Param(key string, value any) Request[R] {
	p.params.Set(key, value)

	return p
}

// Params add params to exisiting params
func (p Request[R]) Params(params Params) Request[R] {
	p.params.Merge(params)
	return p
}

// Data set request data
func (p Request[R]) Data(data io.Reader) Request[R] {
	p.data = data

	return p
}

// Header set request header
func (p Request[R]) Header(key, value string) Request[R] {
	p.headers[key] = value

	return p
}

// Headers set request headers
func (p Request[R]) Headers(headers map[string]string) Request[R] {
	maps.Copy(p.headers, headers)

	return p
}

// NewRequest creates new Request
func NewRequest[R any](api string) Request[R] {
	return Request[R]{
		api:     api,
		method:  http.MethodGet,
		params:  NewParams(),
		headers: make(map[string]string),
	}
}

// Do make request
func (p Request[R]) Do() (R, R, error) {
	api := p.GetURL()

	return MakeJSONRequest[R, R](p.method, api, p.data, p.headers)
}

// GetURL get full URL
func (p Request[R]) GetURL() string {
	return fmt.Sprintf("%s/%s%s", p.api, p.path, p.params.Encode())
}

// BlobResponse binary content response
type BlobResponse struct {
	Content     *bytes.Buffer
	Filename    string
	Size        int64
	ContentType string
}

// Blob make request and returns blob data response
func (p Request[R]) Blob() (r BlobResponse, err error) {
	r.Content = bytes.NewBuffer(nil)
	api := p.GetURL()

	response, err := MakeRequest(p.method, api, p.data, p.headers)

	if err != nil {
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return r, fmt.Errorf("error make request, code: %d status: %s", response.StatusCode, response.Status)
	}

	r.Size, err = r.Content.ReadFrom(response.Body)

	if err != nil {
		return
	}

	r.ContentType = response.Header.Get("Content-Type")

	if _, filename, ok := strings.Cut(response.Header.Get("Content-Disposition"), "filename="); ok {
		r.Filename = strings.ReplaceAll(filename, `"`, "")
	} else {
		r.Filename = "download"
	}

	return
}

// MakeJSONRequest make HTTP request and decode by response by JSON decoder
func MakeJSONRequest[R, E any](method, url string, data io.Reader, headers ...map[string]string) (okResp R, errResp E, err error) {
	response, err := MakeRequest(method, url, data, headers...)

	if err != nil {
		return
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	if response.StatusCode != http.StatusOK {
		if err = decoder.Decode(&errResp); err != nil {
			return
		}

		err = errors.New(response.Status)
	} else {
		err = decoder.Decode(&okResp)
	}

	return
}

// MakeRequest make request and returns *http.Response
//
//	'IMPORTANT: Do not forget close response after reading'
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
