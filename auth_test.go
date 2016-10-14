package shimmie_test

import (
	"testing"

	. "github.com/kusubooru/shimmie"
)

var CookieValueTests = []struct {
	username string
	password string
	ip       string
	out      string
}{
	{"user1", "pass1", "11.11.11.11", Hash(PasswordHash("user1", "pass1") + "11.11.0.0")},
	{"user2", "pass2", "22.22.22.22", Hash(PasswordHash("user2", "pass2") + "22.22.0.0")},
}

func TestCookieValue(t *testing.T) {
	for _, tt := range CookieValueTests {
		phash := Hash(tt.username + tt.password)
		got, want := CookieValue(phash, tt.ip), tt.out
		if got != want {
			t.Errorf("CookieValue(Hash(%q+%q), %q) = %q, want %q\n", tt.username, tt.password, tt.ip, got, want)
		}
	}
}
