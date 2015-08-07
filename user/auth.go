package user

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gplus"
	"net/http"

	"foodtastechess/directory"
)

type Authentication interface {
	CompleteAuthHandler(res http.ResponseWriter, req *http.Request)
	LoginRequired() negroni.HandlerFunc
}

type AuthService struct {
	Config        AuthConfig    `inject:"authConfig"`
	SessionConfig SessionConfig `inject:"sessionConfig"`
	Users         Users         `inject:"users"`

	sessionStore sessions.Store
	provider     goth.Provider
}

type AuthConfig struct {
	GoogleKey    string
	GoogleSecret string
	CallbackUrl  string
	SessionKey   string
}

func NewAuthentication() Authentication {
	return new(AuthService)
}

func (s *AuthService) PreProvide(provider directory.Provider) error {
	err := provider("authConfig", AuthConfig{
		GoogleKey:    "419303763151-c57q5rf3omkr7n3f45a5tfavisovo8jr.apps.googleusercontent.com",
		GoogleSecret: "gDkhFl3VXnVbMBGk7B_MeI2z",
		CallbackUrl:  "http://local.drama9.com:8181/auth/callback",
		SessionKey:   "auth",
	})
	if err != nil {
		return err
	}

	err = provider("sessionConfig", SessionConfig{
		SessionName: "ftc_session",
		Secret:      "secret_123",
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
		log.Debug("Could not get gplus provider: %v", err)
		return err
	}

	return nil
}

func (s *AuthService) LoginRequired() negroni.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		session := GetSession(s.SessionConfig, res, req)
		marshalledAuth, ok := session.Get(s.Config.SessionKey).(string)

		if !ok {
			log.Debug("No session found, creating one.")
			authSession, err := s.provider.BeginAuth(getState(req))
			if err != nil {
				log.Error(fmt.Sprintf("Error creating auth session: %v", err))
				return
			}
			session.Save(s.Config.SessionKey, authSession.Marshal())
			s.loginRedirect(res, req, authSession)
			return
		}

		authSession, err := s.provider.UnmarshalSession(marshalledAuth)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			log.Error("Could not unmarshal auth session: %v", err)
			return
		}

		log.Debug("Auth session: %v", authSession)

		guser, err := s.provider.FetchUser(authSession)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			log.Error("Error fetching user: %v", err)
		}
		log.Debug("User: %v", guser)

		var userId Id
		userId = Id(guser.UserID)
		user := User{
			Id:        userId,
			NickName:  guser.NickName,
			AvatarUrl: guser.AvatarURL,
		}

		context.Set(req, "user", user)

		next(res, req)
	}
}

func (s *AuthService) loginRedirect(res http.ResponseWriter, req *http.Request, authSession goth.Session) {
	url, err := authSession.GetAuthURL()
	if err != nil {
		log.Error(fmt.Sprintf("Could not get Auth URL: %v", err))
	}

	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

func (s *AuthService) CompleteAuthHandler(res http.ResponseWriter, req *http.Request) {
	session := GetSession(s.SessionConfig, res, req)
	marshalledAuth, ok := session.Get(s.Config.SessionKey).(string)
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		log.Error("No session found")
		return
	}

	authSession, err := s.provider.UnmarshalSession(marshalledAuth)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Debug("Could not unmarshal auth session: %v", err)
		return
	}

	_, err = authSession.Authorize(s.provider, req.URL.Query())
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		log.Debug("Could not authorize request, got: %v", err)
	}

	session.Save(s.Config.SessionKey, authSession.Marshal())

	guser, err := s.provider.FetchUser(authSession)
	log.Debug("User: %v", guser)
}

func getState(req *http.Request) string {
	return "state"
}
