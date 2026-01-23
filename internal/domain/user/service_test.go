package user

import (
	"accounter/internal/domain/common"
	"context"
	"math/rand/v2"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type serviceTestSuite struct {
	ctx common.Context
	suite.Suite
	sampleService *UserService
	params        common.RequestParams
}

// TestServiceSuite run serviceTestSuite
func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

// Mock for UserRepository
type mockUserRepostory struct {
	store Users
}

// GetList mock of returning Tasks
func (m *mockUserRepostory) GetList(ctx context.Context, p common.RequestParams) (Users, error) {
	return m.store, nil
}

func (m *mockUserRepostory) GetOne(ctx context.Context, id int64) (User, error) {
	for _, u := range m.store {
		if u.ID == id {
			return u, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (m *mockUserRepostory) Create(ctx context.Context, user *User) error {
	user.ID = rand.Int64N(10)
	m.store = append(m.store, *user)

	return nil
}

func (m *mockUserRepostory) Update(ctx context.Context, user *User) error {
	for i, u := range m.store {
		if u.ID == user.ID {
			m.store[i] = *user
			break
		}
	}

	return nil
}

func (m *mockUserRepostory) Delete(ctx context.Context, id int64) error {
	m.store = slices.DeleteFunc(m.store, func(u User) bool { return u.ID == id })
	return nil
}

func (m *mockUserRepostory) GetByCredentials(ctx context.Context, login, password string) (User, error) {
	for _, u := range m.store {
		if u.Login == login {
			if password != "" && password != u.Password {
				continue
			}

			return u, nil
		}
	}

	return User{}, ErrUserNotFound
}

func (m *mockUserRepostory) WithTx(ctx context.Context, cb func(ctx context.Context) error) error {
	return cb(ctx)
}

func (suite *serviceTestSuite) SetupTest() {
	suite.ctx = common.NewTestContext(context.Background(), 1)
	suite.params = common.NewRequestParams(time.Now())
	repo := &mockUserRepostory{
		store: Users{
			{ID: 1, Name: "Test 1", Login: "Test1"},
			{ID: 2, Name: "Test 1", Login: "Test2"},
			{ID: 3, Name: "Test 1", Login: "Test3"},
		},
	}
	suite.sampleService = NewUserService(repo)
}

func (suite *serviceTestSuite) TestGetUsersList() {
	result, err := suite.sampleService.GetUsersList(suite.ctx, suite.params)

	suite.Equal(err, nil)
	suite.Equal(len(result), 3)
}

func (suite *serviceTestSuite) TestGetUser() {
	user, err := suite.sampleService.GetUser(suite.ctx, 1)

	suite.Nil(err)
	suite.Equal(user.ID, int64(1))

	user, err = suite.sampleService.GetUser(suite.ctx, -1)
	suite.Equal(err, ErrUserNotFound)
}

func (suite *serviceTestSuite) TestSaveUser() {
	var user User

	err := suite.sampleService.SaveUser(suite.ctx, &user)
	suite.Nil(err)
	suite.NotEqual(user.ID, int64(0))

	user = User{ID: 1}
	err = suite.sampleService.SaveUser(suite.ctx, &user)
	suite.Nil(err)
}

func (suite *serviceTestSuite) TestDeleteUser() {
	err := suite.sampleService.DeleteUser(suite.ctx, 1)
	suite.Nil(err)

	err = suite.sampleService.DeleteUser(suite.ctx, -1)
	suite.Equal(err, ErrUserNotFound)
}

func (suite *serviceTestSuite) TestCheckUniqueLogin() {
	exists, err := suite.sampleService.CheckUniqueLogin(suite.ctx, "Test1")
	suite.True(exists)
	suite.Nil(err)

	exists, err = suite.sampleService.CheckUniqueLogin(suite.ctx, "Unknown")
	suite.False(exists)
	suite.Nil(err)
}
