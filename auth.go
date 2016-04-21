package shimmie

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/kusubooru/teian/store"

	"golang.org/x/net/context"
)

// Hash returns the MD5 checksum of a string s as type string.
func Hash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// Auth is a handler wrapper that checks if a user is authenticated to Shimmie.
// It checks for two cookies "shm_user" and "shm_session". The first
// contains the username which is used to query the database and the get user's
// password hash. Then it attempts to recreate the "shm_session" cookie value
// by using the username, user IP and password hash. If the recreated value
// does not match the "shm_session" cookie value then it redirects to
// redirectPath. If redirectURL is empty then "/user_admin/login" is used
// instead which is the default login URL for Shimmie.
func Auth(ctx context.Context, fn func(context.Context, http.ResponseWriter, *http.Request), redirectURL string) http.HandlerFunc {
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
		user, err := store.GetUser(ctx, username)
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
		userIP := getOriginalIP(r)
		sessionCookieValue := CookieValue(passwordHash, userIP)
		if sessionCookieValue != sessionCookie.Value {
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}
		ctx = context.WithValue(ctx, "user", user)
		fn(ctx, w, r)
	}
}

func getOriginalIP(r *http.Request) string {
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
