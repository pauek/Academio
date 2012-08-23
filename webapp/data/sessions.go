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
	Id      string `json:"-"`
	Message string
	Referer string
	Expires time.Time
	User    *User
}

func NewSession(expires time.Time) *Session {
	id := NewUUID()
	log.Printf("New session: %s", id)
	s := &Session{
		Id:      id,
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
	delete(sessions, S.Id)
	mux.Unlock()
}

func (S *Session) PutCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:  cookiename,
		Value: S.Id,
		Path:  "/",
	})
}

func (S *Session) SetUser(user *User) {
	S.User = user
}

func FindSession(id string) *Session {
	s, found := sessions[id]
	if !found {
		return nil
	}
	return s
}

var SessionDuration = time.Hour * 24 * 30 // ~1 month

func GetSession(req *http.Request) *Session {
	cookie, err := req.Cookie(cookiename)
	if err != nil {
		return nil
	}
	return FindSession(cookie.Value)
}

func GetOrCreateSession(req *http.Request) (session *Session) {
	session = GetSession(req)
	if session == nil {
		session = NewSession(time.Now().Add(SessionDuration))
	}
	return
}
