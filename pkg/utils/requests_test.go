package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParams tests Params container
func TestParams(t *testing.T) {
	params := NewParams()
	params.Set("key1", "value1")
	params.Set("key2", 42)

	assert.Equal(t, "value1", params.Get("key1"))
	assert.Equal(t, "42", params.Get("key2"))

	params.Add("key2", "value3")
	assert.Equal(t, "42,value3", params.Get("key2"))

	params.Delete("key1")
	assert.Empty(t, params.Get("key1"))

}

// TestParamsMerge tests merging two Params
func TestParamsMerge(t *testing.T) {
	params1 := NewParams()
	params1.Set("key1", "value1")

	params2 := NewParams()
	params2.Set("key2", "value2")

	params1.Merge(params2)

	assert.Equal(t, "value1", params1.Get("key1"))
	assert.Equal(t, "value2", params1.Get("key2"))
}

// TestParamsEncode tests encoding Params to query string
func TestParamsEncode(t *testing.T) {
	params := NewParams()
	params.Set("key1", "value1")
	params.Set("key2", "value2")

	encoded := params.Encode()
	expected := "?key1=value1&key2=value2"

	assert.Equal(t, expected, encoded)

	paramsEmpty := NewParams()
	encodedEmpty := paramsEmpty.Encode()
	assert.Equal(t, "", encodedEmpty)
}

// TestRequest tests Request builder
func TestReqiest(t *testing.T) {
	params := NewParams()
	params.Set("skip", "0")
	params.Set("limit", "10")

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	data := Data{
		"username": "testuser",
		"email":    "test@example.com",
	}.ToJSON()

	request := NewRequest[Data]("https://example.com").
		Method("POST").
		Path("/test").
		Header("Authorization", "Bearer token").
		Params(params).
		Param("query", "golang").
		Headers(headers).
		Data(data)

	assert.Equal(t, "POST", request.method)
	assert.Equal(t, "https://example.com/test?limit=10&query=golang&skip=0", request.GetURL())
	assert.Equal(t, "Bearer token", request.headers["Authorization"])
	assert.Equal(t, "application/json", request.headers["Content-Type"])
	assert.NotNil(t, request.data)

	request = request.Method("GET").Path("")
	resp, errResp, err := request.Do()

	assert.Nil(t, resp)
	assert.Nil(t, errResp)
	assert.NotNil(t, err)

	blobResp, err := request.Blob()

	assert.Nil(t, err)
	assert.NotNil(t, blobResp.Content)
}
