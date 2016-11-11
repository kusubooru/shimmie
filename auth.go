package shimmie

import (
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type contextKey int

const (
	userContextKey contextKey = iota
)

// Hash returns the MD5 checksum of a string s as type string.
func Hash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// PasswordHash returns the password hash of a username and password the same
// way that shimmie2 does it.
func PasswordHash(username, password string) string {
	hash := md5.Sum([]byte(strings.ToLower(username) + password))
	return fmt.Sprintf("%x", hash)
}

// Auth is a handler wrapper that checks if a user is authenticated to Shimmie.
// It checks for two cookies "shm_user" and "shm_session". The first
// contains the username which is used to query the database and the get user's
// password hash. Then it attempts to recreate the "shm_session" cookie value
// by using the username, user IP and password hash. If the recreated value
// does not match the "shm_session" cookie value then it redirects to
// redirectPath. If redirectURL is empty then "/user_admin/login" is used
// instead which is the default login URL for Shimmie.
func (shim *Shimmie) Auth(fn http.HandlerFunc, redirectURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const defaultLoginURL = "/user_admin/login"
		if redirectURL == "" {
			redirectURL = defaultLoginURL
		}
		usernameCookie, err := r.Cookie("shm_user")
		if err != nil || usernameCookie.Value == "" {
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}
		sessionCookie, err := r.Cookie("shm_session")
		if err != nil {
			log.Print("shimmie: no session cookie")
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}
		username := usernameCookie.Value
		user, err := shim.Store.GetUserByName(username)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("shimmie: user %q does not exist", username)
				http.Redirect(w, r, redirectURL, http.StatusFound)
				return
			}
			msg := fmt.Sprintf("shimmie: could not authenticate: get user %q failed: %v", username, err.Error())
			log.Print(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		passwordHash := user.Pass
		userIP := GetOriginalIP(r)
		sessionCookieValue := CookieValue(passwordHash, userIP)
		if sessionCookieValue != sessionCookie.Value {
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}
		ctx := NewContextWithUser(r.Context(), user)
		fn(w, r.WithContext(ctx))
	}
}

// FromContextGetUser gets User from context. If User does not exist in context,
// nil and false are returned instead.
func FromContextGetUser(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userContextKey).(*User)
	return user, ok
}

// NewContextWithUser adds user to context.
func NewContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// GetOriginalIP gets the original IP of the HTTP for the case of being behind
// a proxy. It searches for the X-Forwarded-For header.
func GetOriginalIP(r *http.Request) string {
	x := r.Header.Get("X-Forwarded-For")
	if x != "" && strings.Contains(r.RemoteAddr, "127.0.0.1") {
		// format is comma separated
		return strings.Split(x, ",")[0]
	}
	// it also contains the port
	return strings.Split(r.RemoteAddr, ":")[0]
}

// CookieValue recreates the Shimmie session cookie value based on the user
// password hash and the user IP.
//
// Shimmie creates a cookie "shm_session" containing an md5 digest value of the
// user password hash concatenated with the user IP masked with the 255.255.0.0
// mask. That's essentially:
//
//   md5(password_hash + masked_ip)
//
func CookieValue(passwordHash, userIP string) string {
	addr := net.ParseIP(strings.Split(userIP, ":")[0])
	mask := net.IPv4Mask(255, 255, 0, 0)
	addr = addr.Mask(mask)
	sessionHash := md5.Sum([]byte(fmt.Sprintf("%s%s", passwordHash, addr)))
	return fmt.Sprintf("%x", sessionHash)
}

const loginMemory = 365

// SetCookie creates a cookie on path "/" with 1 year expiration and other
// flags set to false mimicking the cookies that Shimmie creates.
func SetCookie(w http.ResponseWriter, name, value string) {
	expires := time.Now().Add(time.Second * 60 * 60 * 24 * loginMemory)
	c := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expires,
		Path:    "/",
	}
	http.SetCookie(w, &c)
}
