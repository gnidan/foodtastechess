package auth

import (
	"errors"
	"fmt"
	"github.com/gorilla/context"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gplus"
	"net/http"
	"net/url"

	"foodtastechess/directory"
	"foodtastechess/logger"
	sess "foodtastechess/server/session"
	"foodtastechess/user"
)

var log = logger.Log("auth")

type Authentication interface {
	LoginRequired(res http.ResponseWriter, req *http.Request, next http.HandlerFunc)
}

type authService struct {
	Config        AuthConfig         `inject:"authConfig"`
	SessionConfig sess.SessionConfig `inject:"sessionConfig"`
	Users         user.Users         `inject:"users"`

	provider goth.Provider
}

func New() Authentication {
	return new(authService)
}

// PreProvide is just for creating a fake auth config at this point
func (s *authService) PreProvide(provider directory.Provider) error {
	err := provider("authConfig", AuthConfig{
		GoogleKey:    "419303763151-c57q5rf3omkr7n3f45a5tfavisovo8jr.apps.googleusercontent.com",
		GoogleSecret: "gDkhFl3VXnVbMBGk7B_MeI2z",
		CallbackUrl:  "http://local.drama9.com:8181/auth/callback",
		SessionKey:   "auth",
	})

	return err
}

// PostPopulate sets up the oauth provider
func (s *authService) PostPopulate() error {
	goth.UseProviders(gplus.New(
		s.Config.GoogleKey,
		s.Config.GoogleSecret,
		s.Config.CallbackUrl,
	))

	var err error
	s.provider, err = goth.GetProvider("gplus")
	if err != nil {
		log.Error(fmt.Sprintf("Could not get gplus provider: %v", err))
		return err
	}

	return nil
}

// LoginRequired is a middleware function that catches certain paths and blocks
// the next handler from occuring unless the user is logged in.
//
// Authentication-specific routes that will be handled by LoginRequired:
//
// - /auth/login
//
//	 Starts a new auth session and redirects to the OAuth provider
//
// - /auth/callback
//
//   Catches the redirect back from the provider and saves the auth session
//
// - /auth/me
//
//	 Returns some information about the user or responds `401 Unauthorized`
//
func (s *authService) LoginRequired(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	session := sess.GetSession(s.SessionConfig, res, req)

	switch req.URL.Path {
	case "/auth/login":
		s.beginAuth(res, req, session)
	case "/auth/callback":
		s.completeAuth(res, req, session)
	case "/auth/me":
		s.authInfo(res, req, session)
	default:
		u, valid := s.validCredentials(session)
		if !valid {
			log.Info("Not Logged In, Redirecting")

			loginPath := "/auth/login"

			params := url.Values{}
			params.Set("redirect", url.QueryEscape(req.URL.String()))

			redirectUrl := fmt.Sprintf("%s?%s", loginPath, params.Encode())

			http.Redirect(res, req, redirectUrl, http.StatusTemporaryRedirect)
			return
		}

		context.Set(req, ContextKey, u)
		next(res, req)
	}
}

func (s *authService) validCredentials(session sess.Session) (user.User, bool) {
	u := user.User{}

	authSession, err := s.loadAuthSession(session)
	if err != nil {
		log.Info(fmt.Sprintf("Valid Auth Session not found: %v", err))
		return u, false
	}

	guser, err := s.provider.FetchUser(authSession)
	if err != nil {
		log.Error(fmt.Sprintf("Error fetching user: %v", err))
		return u, false
	}

	if guser.RawData["error"] != nil {
		log.Info("User Session Expired")
		return u, false
	}

	u = user.User{
		AuthIdentifier: guser.UserID,
		Name:           guser.NickName,
		AvatarUrl:      guser.AvatarURL,
	}

	return u, true
}

func (s *authService) beginAuth(res http.ResponseWriter, req *http.Request, session sess.Session) {
	authSession, err := s.provider.BeginAuth(getState(req))
	if err != nil {
		log.Error(fmt.Sprintf("Error creating auth session: %v", err))
		return
	}
	session.Save(s.Config.SessionKey, authSession.Marshal())

	url, err := authSession.GetAuthURL()
	if err != nil {
		log.Error(fmt.Sprintf("Could not get Auth URL: %v", err))
	}

	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

func (s *authService) completeAuth(res http.ResponseWriter, req *http.Request, session sess.Session) {
	authSession, err := s.loadAuthSession(session)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		log.Error(fmt.Sprintf("Could not load auth session: %v", err))
		return
	}

	_, err = authSession.Authorize(s.provider, req.URL.Query())
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		log.Info("Could not authorize request, got: %v", err)
	}

	s.saveAuthSession(session, authSession)

	redirectUrl, _ := url.QueryUnescape(req.URL.Query().Get("state"))
	http.Redirect(res, req, redirectUrl, http.StatusTemporaryRedirect)
}

func (s *authService) authInfo(res http.ResponseWriter, req *http.Request, session sess.Session) {
	u, valid := s.validCredentials(session)
	if !valid {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	res.Write([]byte(fmt.Sprintf("%v", u.AuthIdentifier)))
}

func (s *authService) saveAuthSession(session sess.Session, authSession goth.Session) {
	session.Save(s.Config.SessionKey, authSession.Marshal())
}

func (s *authService) loadAuthSession(session sess.Session) (goth.Session, error) {
	marshalledAuth, ok := session.Get(s.Config.SessionKey).(string)
	if !ok {
		return nil, errors.New("No auth session found")
	}

	authSession, err := s.provider.UnmarshalSession(marshalledAuth)
	if err != nil {
		return nil, err
	}

	return authSession, nil
}

func getState(req *http.Request) string {
	redirect := req.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/auth/me"
	}
	return redirect
}
