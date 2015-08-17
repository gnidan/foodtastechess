package auth

import (
	"errors"
	"fmt"
	"github.com/gorilla/context"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gplus"
	"net/http"
	"net/url"

	"foodtastechess/config"
	"foodtastechess/logger"
	sess "foodtastechess/server/session"
	"foodtastechess/user"
)

var log = logger.Log("auth")

type Authentication interface {
	LoginRequired(res http.ResponseWriter, req *http.Request, next http.HandlerFunc)
}

type authService struct {
	Config        config.AuthConfig    `inject:"authConfig"`
	SessionConfig config.SessionConfig `inject:"sessionConfig"`
	Users         user.Users           `inject:"users"`

	provider goth.Provider
}

func New() Authentication {
	return new(authService)
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
	authSession, err := s.loadAuthSession(session)
	if err != nil {
		log.Info(fmt.Sprintf("Valid Auth Session not found: %v", err))
		return user.User{}, false
	}

	return s.getUser(authSession)
}

func (s *authService) getUser(authSession goth.Session) (user.User, bool) {
	guser, err := s.provider.FetchUser(authSession)
	if err != nil {
		log.Error(fmt.Sprintf("Error fetching user: %v", err))
		return user.User{}, false
	}

	if guser.RawData["error"] != nil {
		log.Info("User Session Expired")
		return user.User{}, false
	}

	u, found := s.Users.GetByAuthId(guser.UserID)
	if found {
		u.Name = guser.NickName
		u.AvatarUrl = guser.AvatarURL
		u.AccessToken = guser.AccessToken
		u.AccessTokenSecret = guser.AccessTokenSecret
	} else {
		u = user.User{
			Name:              guser.NickName,
			AvatarUrl:         guser.AvatarURL,
			AuthIdentifier:    guser.UserID,
			Uuid:              user.NewId(),
			AccessToken:       guser.AccessToken,
			AccessTokenSecret: guser.AccessTokenSecret,
		}
	}

	s.Users.Save(&u)

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

	// Triggers a save
	s.getUser(authSession)

	redirectUrl, _ := url.QueryUnescape(req.URL.Query().Get("state"))
	http.Redirect(res, req, redirectUrl, http.StatusTemporaryRedirect)
}

func (s *authService) authInfo(res http.ResponseWriter, req *http.Request, session sess.Session) {
	u, valid := s.validCredentials(session)
	if !valid {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	res.Write([]byte(fmt.Sprintf("%v", u.Uuid)))
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
