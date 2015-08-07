package server

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type Session interface {
	Get(key string) interface{}
	Save(key string, value interface{})
	Delete(key string)
}

type CookieSession struct {
	config SessionConfig
	store  sessions.Store
	res    http.ResponseWriter
	req    *http.Request
}

func GetSession(config SessionConfig, res http.ResponseWriter, req *http.Request) Session {
	session := new(CookieSession)
	session.config = config
	session.store = sessions.NewCookieStore([]byte(config.Secret))
	session.res = res
	session.req = req

	return session
}

func (s *CookieSession) Get(key string) interface{} {
	gorillaSession, _ := s.store.Get(s.req, s.config.SessionName)
	return gorillaSession.Values[key]
}

func (s *CookieSession) Save(key string, value interface{}) {
	gorillaSession, _ := s.store.Get(s.req, s.config.SessionName)
	gorillaSession.Values[key] = value
	gorillaSession.Save(s.req, s.res)
}

func (s *CookieSession) Delete(key string) {
	gorillaSession, _ := s.store.Get(s.req, s.config.SessionName)
	delete(gorillaSession.Values, key)
	gorillaSession.Save(s.req, s.res)
}
