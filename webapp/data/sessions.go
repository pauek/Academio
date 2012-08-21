package data

import (
	"log"
	"net/http"
	"sync"
	"time"
)

const cookiename = "session-id"

var (
	sessions = make(map[string]*Session)
	mux      sync.Mutex
)

type Session struct {
	id      string `json:"-"`
	rev     string `json:"-"`
	Message string
	Referer string
	Expires time.Time
	User    interface{}
}

func NewSession(expires time.Time) *Session {
	id := NewUUID()
	log.Printf("New session: %s", id)
	s := &Session{
		id:      id,
		User:    nil,
		Expires: expires,
	}
	mux.Lock()
	sessions[id] = s
	mux.Unlock()
	return s
}

func (S *Session) Delete() {
	mux.Lock()
	delete(sessions, S.id)
	mux.Unlock()
}

func (S *Session) SetCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:  cookiename,
		Value: S.id,
		Path:  "/",
	})
}

func (S *Session) SetUser(user interface{}) {
	S.User = user
}

func GetSession(id string) *Session {
	s, found := sessions[id]
	if !found {
		return nil
	}
	return s
}

func GetSessionFromRequest(req *http.Request) *Session {
	cookie, err := req.Cookie(cookiename)
	if err != nil {
		return nil
	}
	return GetSession(cookie.Value)
}

