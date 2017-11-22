package apiuser

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// START OMIT
func TestHandler(t *testing.T) {
	assert := assert.New(t) // testify

	// make test w and r
	r, err := http.NewRequest("GET", "http://users/listusers", nil)
	fmt.Println(r.URL.Path)
	fmt.Println(r.URL.Scheme)
	r.URL.Path = "http://users/listusers"
	assert.NoError(err)
	w := httptest.NewRecorder()
	var dbUser = make(data.DBtype)
	us := NewUserService(dbUser)
	u := NewUserHandler(us)
	h := NewHandler(u)

	//ctx := context.WithValue(r.Context(), MyKey, &Claims{"chidi", "", jwt.StandardClaims{}})
	_ = u.Userservice.AddUser(&data.User{Username: "chidi", Password: "cc", Level: "Admin"})
	sr, err := u.Userservice.GetUser("chidi")
	fmt.Println("User = ", sr)
	// call handler
	fmt.Println(r.URL.Path)
	h.ServeHTTP(w, r)
	// make assertions
	//assert.Equal(w.Body.String(), "Hello 9ja", "body")
	assert.Equal(w.Code, http.StatusOK, "status code")
}

// END OMIT

// func TestMain(t *testing.T) {
// 	var tests []testing.InternalTest
// 	tests = append(tests, testing.InternalTest{Name: "TestHandler", F: TestHandler})
// 	testing.Main(func(pat, str string) (bool, error) { return true, nil }, tests, nil, nil)
// }

// package main

// import (
// 	"bytes"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// )

// func TestMainHandler(t *testing.T) {
// 	rootRequest, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Fatal("Root request error: %s", err)
// 	}

// 	cases := []struct {
// 		w                    *httptest.ResponseRecorder
// 		r                    *http.Request
// 		accessTokenHeader    string
// 		expectedResponseCode int
// 		expectedResponseBody []byte
// 	}{
// 		{
// 			w:                    httptest.NewRecorder(),
// 			r:                    rootRequest,
// 			accessTokenHeader:    "magic",
// 			expectedResponseCode: http.StatusOK,
// 			expectedResponseBody: []byte("You have some magic in you\n"),
// 		},
// 		{
// 			w:                    httptest.NewRecorder(),
// 			r:                    rootRequest,
// 			accessTokenHeader:    "",
// 			expectedResponseCode: http.StatusForbidden,
// 			expectedResponseBody: []byte("You don't have enough magic in you\n"),
// 		},
// 	}

// 	for _, c := range cases {
// 		c.r.Header.Set("X-Access-Token", c.accessTokenHeader)

// 		mainHandler(c.w, c.r)

// 		if c.expectedResponseCode != c.w.Code {
// 			t.Errorf("Status Code didn't match:\n\t%q\n\t%q", c.expectedResponseCode, c.w.Code)
// 		}

// 		if !bytes.Equal(c.expectedResponseBody, c.w.Body.Bytes()) {
// 			t.Errorf("Body didn't match:\n\t%q\n\t%q", string(c.expectedResponseBody), c.w.Body.String())
// 		}
// 	}
// }
