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

type MockUsers struct {
	mock.Mock
}

func (m *MockUsers) Get(id Id) (User, bool) {
	return User{}, false
}

func (m *MockUsers) Save(user User) error {
	return nil
}

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

func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
