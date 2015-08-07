package auth

import (
	"fmt"
	"github.com/gorilla/context"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gplus"
	"net/http"
	"net/url"

	"foodtastechess/directory"
	"foodtastechess/logger"
	"foodtastechess/server/session"
	"foodtastechess/user"
)

var log = logger.Log("auth")

type Authentication interface {
	CompleteAuthHandler(res http.ResponseWriter, req *http.Request)
	LoginRequired(res http.ResponseWriter, req *http.Request, next http.HandlerFunc)
}

type AuthService struct {
	Config        AuthConfig            `inject:"authConfig"`
	SessionConfig session.SessionConfig `inject:"sessionConfig"`
	Users         user.Users            `inject:"users"`

	provider goth.Provider
}

func New() Authentication {
	return new(AuthService)
}

func (s *AuthService) PreProvide(provider directory.Provider) error {
	err := provider("authConfig", AuthConfig{
		GoogleKey:    "419303763151-c57q5rf3omkr7n3f45a5tfavisovo8jr.apps.googleusercontent.com",
		GoogleSecret: "gDkhFl3VXnVbMBGk7B_MeI2z",
		CallbackUrl:  "http://local.drama9.com:8181/auth/callback",
		SessionKey:   "auth",
	})

	return err
}

func (s *AuthService) PostPopulate() error {
	goth.UseProviders(gplus.New(
		s.Config.GoogleKey,
		s.Config.GoogleSecret,
		s.Config.CallbackUrl,
	))

	var err error
	s.provider, err = goth.GetProvider("gplus")
	if err != nil {
		log.Error("Could not get gplus provider: %v", err)
		return err
	}

	return nil
}

func (s *AuthService) LoginRequired(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if req.URL.Path == "/auth/callback" {
		s.CompleteAuthHandler(res, req)
		return
	}

	sess := session.GetSession(s.SessionConfig, res, req)
	marshalledAuth, ok := sess.Get(s.Config.SessionKey).(string)

	if !ok {
		log.Info("No session found, creating one.")
		s.startAuth(res, req, sess)
		return
	}

	authSession, err := s.provider.UnmarshalSession(marshalledAuth)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Error(fmt.Sprintf("Could not unmarshal auth session: %v", err))
		return
	}

	guser, err := s.provider.FetchUser(authSession)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Error(fmt.Sprintf("Error fetching user: %v", err))
	}

	if guser.RawData["error"] != nil {
		log.Info("No access token found, starting auth over")
		s.startAuth(res, req, sess)
		return
	}

	u := user.User{
		Id:        user.Id(guser.UserID),
		NickName:  guser.NickName,
		AvatarUrl: guser.AvatarURL,
	}

	context.Set(req, ContextKey, u)

	next(res, req)
}

func (s *AuthService) startAuth(res http.ResponseWriter, req *http.Request, sess session.Session) {
	authSession, err := s.provider.BeginAuth(getState(req))
	if err != nil {
		log.Error(fmt.Sprintf("Error creating auth session: %v", err))
		return
	}
	sess.Save(s.Config.SessionKey, authSession.Marshal())
	s.loginRedirect(res, req, authSession)
}

func (s *AuthService) loginRedirect(res http.ResponseWriter, req *http.Request, authSession goth.Session) {
	url, err := authSession.GetAuthURL()
	if err != nil {
		log.Error(fmt.Sprintf("Could not get Auth URL: %v", err))
	}

	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

func (s *AuthService) CompleteAuthHandler(res http.ResponseWriter, req *http.Request) {
	sess := session.GetSession(s.SessionConfig, res, req)
	marshalledAuth, ok := sess.Get(s.Config.SessionKey).(string)
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		log.Error("No session found")
		return
	}

	authSession, err := s.provider.UnmarshalSession(marshalledAuth)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Error("Could not unmarshal auth session: %v", err)
		return
	}

	_, err = authSession.Authorize(s.provider, req.URL.Query())
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		log.Info("Could not authorize request, got: %v", err)
	}

	sess.Save(s.Config.SessionKey, authSession.Marshal())

	//guser, _ := s.provider.FetchUser(authSession)

	redirectUrl, _ := url.QueryUnescape(req.URL.Query().Get("state"))
	http.Redirect(res, req, redirectUrl, http.StatusTemporaryRedirect)
}

func getState(req *http.Request) string {
	state := url.QueryEscape(req.URL.String())
	return state
}
