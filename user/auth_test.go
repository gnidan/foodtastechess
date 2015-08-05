package user

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"

	"foodtastechess/directory"
)

// MockUsers mocks the users service for us
type MockUsers struct {
	mock.Mock
}

func (m *MockUsers) Get(id Id) (User, bool) {
	args := m.Called(id)
	return args.Get(0).(User), args.Bool(1)
}

func (m *MockUsers) Save(user User) error {
	args := m.Called(user)
	return args.Error(0)
}

// AuthTestSuite provides a setup by which to test the behavior
// of the AuthService. It mocks the users service and uses
// a fake AuthConfig.
type AuthTestSuite struct {
	suite.Suite

	auth Authentication
}

func (suite *AuthTestSuite) SetupTest() {
	var (
		d directory.Directory

		users MockUsers
	)

	auth := NewAuthentication()
	authConfig := AuthConfig{
		GoogleKey:    "key",
		GoogleSecret: "secret",
		CallbackUrl:  "http://callback/",
	}

	d = directory.New()
	d.AddService("auth", auth)
	d.AddService("authConfig", authConfig)
	d.AddService("users", &users)

	if err := d.Start(); err != nil {
		log.Fatalf("Could not start directory (%v)", err)
	}

	suite.auth = auth
}

func (suite *AuthTestSuite) TestBeginAuth() {
	handle := suite.auth.Handler()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handle.ServeHTTP(w, req)
	assert := assert.New(suite.T())
	assert.Equal(
		w.Code, http.StatusTemporaryRedirect,
		fmt.Sprintf("Got response: '%v'", w.Body),
	)
}

// TestAuth kicks off the AuthTestSuite
func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
