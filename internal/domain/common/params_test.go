package common

import (
	"accounter/pkg/utils"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// RequestParams testing suite
type testRequestParamsSuite struct {
	suite.Suite
	dateStart    time.Time
	dateEnd      time.Time
	sampleParams RequestParams
}

// TestRequestParamsSuite run Task testing
func TestRequestParamsSuite(t *testing.T) {
	suite.Run(t, new(testRequestParamsSuite))
}

// SetupTest testing preparations
func (suite *testRequestParamsSuite) SetupTest() {
	suite.dateStart = time.Date(2026, 1, 1, 0, 0, 0, 0, time.Local)
	suite.dateEnd = time.Date(2026, 1, 1, 23, 59, 59, 0, time.Local)

	suite.sampleParams = RequestParams{
		DateStart: suite.dateStart.Unix(),
		DateEnd:   suite.dateEnd.Unix(),
		Timezone:  "Local",
		OrderBy:   "date",
		Users:     []int64{1, 2},
	}
}

func (suite *testRequestParamsSuite) TestNewParams() {
	params := NewRequestParams(suite.dateStart)

	suite.Equal(params.DateStart, suite.dateStart.Unix(), "Incorrect date start")
	suite.Equal(params.DateEnd, suite.dateEnd.Unix(), "Incorrect date start")

	suite.Equal(params.OrderBy, "date", "Incorrect order key")
	suite.Equal(params.OrderDesc, false, "Incorrect order desc")
	suite.Equal(params.Timezone, time.Local.String(), "Incorrect timezone")
	suite.Empty(params.Status, "Incorrect status")
	suite.Empty(params.Skip, "Incorrect skip")
	suite.Empty(params.Limit, "Incorrect limit")
	suite.Empty(params.Users, "Incorrect users")
}

func (suite *testRequestParamsSuite) TestEncode() {
	encoded := suite.sampleParams.Encode()

	tests := []struct {
		name     string
		expected string
		key      string
	}{
		{name: "Date start existance",
			expected: fmt.Sprintf("%d", suite.sampleParams.DateStart),
			key:      "date_start",
		},
		{
			name:     "Date end existance",
			expected: fmt.Sprintf("%d", suite.sampleParams.DateEnd),
			key:      "date_end",
		},
		{
			name:     "Timezone existance",
			expected: suite.sampleParams.Timezone,
			key:      "timezone",
		},
		{
			name:     "Status existance",
			expected: suite.sampleParams.Status,
			key:      "status",
		},
		{
			name:     "Users existance",
			expected: utils.Stringify(suite.sampleParams.Users...),
			key:      "users",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.Equal(tt.expected, encoded.Get(tt.key))
		})
	}
}
