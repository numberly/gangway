package session

import (
	"crypto/sha256"
	"net/http"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
)

// Session defines a Gangway session
type Session struct {
	Session *CustomCookieStore
}

// New inits a Session with CookieStore
func New(sessionSecurityKey string, sessionSalt string) *Session {
	return &Session{
		Session: NewCustomCookieStore(generateSessionKeys(sessionSecurityKey, sessionSalt)),
	}
}

// generateSessionKeys creates a signed encryption key for the cookie store
func generateSessionKeys(sessionSecurityKey string, salt string) ([]byte, []byte) {
	// Take the configured security key and generate 96 bytes of data. This is
	// used as the signing and encryption keys for the cookie store.  For details
	// on the PBKDF2 function: https://en.wikipedia.org/wiki/PBKDF2
	b := pbkdf2.Key(
		[]byte(sessionSecurityKey),
		[]byte(salt),
		4096, 96, sha256.New)

	return b[0:64], b[64:96]
}

// Cleanup removes the current session from the store
func (s *Session) Cleanup(w http.ResponseWriter, r *http.Request, name string) {
	session, err := s.Session.Get(r, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		log.Errorf("failed to save session: %v", err)
	}
}
