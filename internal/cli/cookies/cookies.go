package cookies

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/docker/distribution/uuid"
	"net/http"
	"time"
)

var secretkey = []byte("CHANGE_ME")

// TODO: add to main func
func SetSecret(secret []byte) {
	secretkey = secret
}

func SetCookie(readyCookie ...string) *http.Cookie {
	fid := uuid.Generate()
	cookie := new(http.Cookie)
	cookie.Name = "session"
	if len(readyCookie) > 0 {
		cookie.Value = readyCookie[0]
	} else {
		cookie.Value = NewCookie(fid.String()).Value + fid.String()
	}
	cookie.Expires = time.Now().Add(24 * time.Hour * 365)
	return cookie
}

func NewCookie(username string) *http.Cookie {
	h := hmac.New(sha256.New, secretkey)
	src := []byte(username)
	h.Write(src)

	value := hex.EncodeToString(h.Sum(nil)) + "-" + hex.EncodeToString(src)
	cookie := &http.Cookie{
		Name:       "session",
		Value:      value,
		Path:       "",
		Domain:     "localhost",
		Expires:    time.Time{},
		RawExpires: "",
		MaxAge:     3600,
		Secure:     false,
		HttpOnly:   true,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}

	return cookie
}
