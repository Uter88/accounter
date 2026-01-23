package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

// testCurrentUserSuite test suite for User
type testCurrentUserSuite struct {
	suite.Suite
	sampleUser CurrentUser
}

// TestCurrentUserSuitee run testCurrentUserSuite
func TestCurrentUserSuitee(t *testing.T) {
	suite.Run(t, new(testCurrentUserSuite))
}

// SetupTest test preparations
func (suite *testCurrentUserSuite) SetupTest() {
	suite.sampleUser = CurrentUser{}
}

func (suite *testCurrentUserSuite) TestSetToken() {
	suite.sampleUser.SetToken("access_token", "refesh_token")

	suite.Equal(suite.sampleUser.Tokens.AccessToken, "access_token")
	suite.Equal(suite.sampleUser.Tokens.RefreshToken, "refesh_token")
	suite.True(suite.sampleUser.IsAuthorized)
}

func (suite *testCurrentUserSuite) TestSetAuthorized() {
	suite.sampleUser.SetAuthorized(true)
	suite.True(suite.sampleUser.IsAuthorized)
}

func (suite *testCurrentUserSuite) TestSetContext() {
	ctx := context.Background()
	user := suite.sampleUser.WithContext(ctx)
	suite.Equal(user.Context, ctx)
}
