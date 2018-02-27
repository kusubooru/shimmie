package store_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/kusubooru/shimmie"
)

func TestGetPMs(t *testing.T) {
	shim, schema := setup(t)
	defer teardown(t, shim, schema)

	var users = []*shimmie.User{
		&shimmie.User{Name: "bob", Pass: "bob123"},
		&shimmie.User{Name: "ann", Pass: "ann123"},
		&shimmie.User{Name: "zoe", Pass: "zoe123"},
	}
	for _, u := range users {
		if err := shim.CreateUser(u); err != nil {
			t.Fatalf("CreateUser(%q) returned err: %v", u, err)
		}
	}

	var pms = []*shimmie.PM{
		&shimmie.PM{FromID: 1, ToID: 2, Subject: "bob greets ann", SentDate: time.Now(), IsRead: true},
		&shimmie.PM{FromID: 2, ToID: 1, Subject: "ann greets bob", SentDate: time.Now(), IsRead: true},
		&shimmie.PM{FromID: 2, ToID: 3, Subject: "ann greets zoe", SentDate: time.Now(), IsRead: true},
		&shimmie.PM{FromID: 1, ToID: 3, Subject: "bob greets zoe", SentDate: time.Now(), IsRead: true},

		&shimmie.PM{FromID: 1, ToID: 2, Subject: "bob texts ann (Unread)", SentDate: time.Now()},
		&shimmie.PM{FromID: 2, ToID: 1, Subject: "ann texts bob (Unread)", SentDate: time.Now()},
		&shimmie.PM{FromID: 2, ToID: 3, Subject: "ann texts zoe (Unread)", SentDate: time.Now()},
		&shimmie.PM{FromID: 1, ToID: 3, Subject: "bob texts zoe (Unread)", SentDate: time.Now()},
	}
	for _, pm := range pms {
		if err := shim.CreatePM(pm); err != nil {
			t.Fatalf("error creating pm %#v: %v", pm, err)
		}
	}

	var tests = []struct {
		from string
		to   string
		r    shimmie.PMChoice
		want int
	}{
		{want: 8},
		{want: 4, r: shimmie.PMRead},
		{want: 4, r: shimmie.PMUnread},

		{from: "bob", to: "ann", want: 1, r: shimmie.PMRead},
		{from: "ann", to: "bob", want: 1, r: shimmie.PMRead},
		{from: "bob", want: 2, r: shimmie.PMRead},
		{from: "ann", want: 2, r: shimmie.PMRead},
		{to: "zoe", want: 2, r: shimmie.PMRead},

		{from: "bob", to: "ann", want: 2, r: shimmie.PMAny},
	}

	for i, tt := range tests {
		pms, err := shim.GetPMs(tt.from, tt.to, tt.r)
		if err != nil {
			t.Fatalf("%d: getting pms from %q, to %q, choice %d returned err: %v", i, tt.from, tt.to, tt.r, err)
		}
		if got, want := len(pms), tt.want; got != want {
			t.Errorf("%d: getting pms from %q, to %q, choice %d -> len(pms) = %d, want %d", i, tt.from, tt.to, tt.r, got, want)
			data, _ := json.Marshal(pms)
			fmt.Println(string(data))
		}
	}
}

func TestCreatePM(t *testing.T) {
	shim, schema := setup(t)
	defer teardown(t, shim, schema)

	err := shim.CreatePM(nil)
	if err == nil {
		t.Errorf("creating nil PM should return error")
	}
}
