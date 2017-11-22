//2nd in this package
package apiuser

import (
	"Game/tools"
	"context"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Key int

var MyKey Key = 0

type Session struct {
	db              interface{}
	Userservicesess UserServiceSess
}

func NewSession(uDB DBType) *Session {
	s := &Session{
		db: uDB,
	}
	s.Userservicesess.session = s
	return s
}

// JWT schema of the data it will store
type Claims struct {
	Username    string `json:"username"`
	RedirectURL string
	Level       string
	jwt.StandardClaims
}

// create a JWT and put in the clients cookie
func (s *Session) setToken(w http.ResponseWriter, r *http.Request, username, redirectURL, level string) {
	expireCookie := time.Now().Add(time.Minute * 20)
	expireToken := time.Now().Add(time.Minute * 20).Unix()
	claims := &Claims{
		username,
		redirectURL,
		level,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:8080",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte("secret"))
	cookie := http.Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true}
	http.SetCookie(w, &cookie)
	time.Sleep(300 * time.Millisecond)
	// redirect
	http.Redirect(w, r, claims.RedirectURL, http.StatusSeeOther)
	return
}

// middleware to protect private pages
func (s *Session) Validate(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Auth")
		if err != nil {
			ctx := context.WithValue(r.Context(), MyKey, &Claims{"anonymous", "", "", jwt.StandardClaims{}})
			handler(w, r.WithContext(ctx))
			return
		}
		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), MyKey, *claims)
			handler(w, r.WithContext(ctx))
		} else {
			http.NotFound(w, r)
			return
		}
	})
}

// deletes the cookie
func (s *Session) logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("Auth")
	if err != nil {
		return
	}
	deleteCookie := http.Cookie{Name: "Auth", Value: "none", Expires: time.Now()}
	http.SetCookie(w, &deleteCookie)
	return
}

func (s *Session) UserFromRequest(r *http.Request) (*User, error) {
	claims, ok := r.Context().Value(MyKey).(Claims)
	if !ok || claims.Username == "anomnymous" {
		return nil, tools.ErrUnauthorized
	}
	user, err := s.Userservicesess.GetUser(claims.Username)
	if err != nil {
		return nil, tools.ErrUnauthorized
	}
	return user, nil
}
