package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type errorTestSuite struct {
	suite.Suite
}

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(errorTestSuite))
}

func (suite *errorTestSuite) TestErrorAndWithError() {
	err := &DomainError{Code: "error", Message: "test error"}
	suite.Equal(err.Error(), "error: test error")

	err = err.WithErr(errors.New("empty"))
	suite.Equal(err.Error(), "error: test error (caused by: empty)")
}

func (suite *errorTestSuite) TestIsDomainError() {
	err := &DomainError{Code: "error"}
	suite.True(IsDomainError(err, "error"))

	suite.False(IsDomainError(errors.New("error"), "error"))
}
