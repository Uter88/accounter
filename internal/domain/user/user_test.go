package user

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// testUserSuite test suite for User
type testUserSuite struct {
	suite.Suite
	sampleUser User
}

// TestUserSuite run testUserSuite
func TestUserSuite(t *testing.T) {
	suite.Run(t, new(testUserSuite))
}

// SetupTest test preparations
func (s *testUserSuite) SetupTest() {
	s.sampleUser = User{
		ID:           1,
		Name:         "Test",
		Surname:      "Test",
		Patronymic:   "Test",
		Login:        "test@gmail.com",
		Password:     "Test",
		PricePerHour: 1500,
		MoneyEarned:  355000,
	}
}

func (suite *testUserSuite) TestGetIDAndType() {
	task := suite.sampleUser

	suite.Equal(int64(1), task.GetID())
	suite.Equal(UserType, task.GetType())
}

// TestFormatMethods test format methods
func (suite *testUserSuite) TestFormatMethods() {
	suite.Equal("Test T. T.", suite.sampleUser.GetLabel())
	suite.Equal(float32(355000.0), suite.sampleUser.GetMoneyEarned())
	suite.Equal("355 000.00", suite.sampleUser.FormatMoneyEarned())
}

// TestIsValid test User validation
func (suite *testUserSuite) TestIsValid() {
	user := suite.sampleUser

	suite.Equal(true, user.IsValid(false))

	user.PricePerHour = 0
	suite.Equal(false, user.IsValid(false))

	user.Name = ""
	suite.Equal(false, user.IsValid(false))

	user.Login = "test@gmail.com"
	user.Password = ""
	suite.Equal(false, user.IsValid(true))

	user.Login = "bad email"
	suite.Equal(false, user.IsValid(true))
}
