package user

import (
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/gplus"
	"net/http"
)

type Authentication interface {
	Handler() http.Handler
}

type AuthService struct {
	Config AuthConfig `inject:"authConfig"`
	Users  Users      `inject:"users"`
}

type AuthConfig struct {
	GoogleKey    string
	GoogleSecret string
	CallbackUrl  string
}

func NewAuthentication() Authentication {
	return new(AuthService)
}

func (s *AuthService) PostPopulate() error {
	goth.UseProviders(
		gplus.New(
			s.Config.GoogleKey,
			s.Config.GoogleSecret,
			s.Config.CallbackUrl,
		),
	)

	return nil
}

func (s *AuthService) Handler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", s.BeginAuth)
	r.HandleFunc("/callback", s.CompleteAuth)
	return r
}

func (s *AuthService) BeginAuth(res http.ResponseWriter, req *http.Request) {
	log.Debug("%v", req.URL.Query)
	gothic.GetProviderName = getProviderName
	gothic.BeginAuthHandler(res, req)
}

func (s *AuthService) CompleteAuth(res http.ResponseWriter, req *http.Request) {
	guser, err := gothic.CompleteUserAuth(res, req)

	if err != nil {
		log.Error("Something went wrong")
		return
	}

	log.Debug("%v", guser)
}

func getProviderName(req *http.Request) (string, error) {
	return "gplus", nil
}
