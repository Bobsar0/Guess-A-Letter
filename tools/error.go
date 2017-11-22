package tools

import "errors"

// General errors.
var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInternal     = errors.New("internal error")
)

// User and Movie errors.
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrMovieNotFound    = errors.New("movie not found")
	ErrGameNotFound     = errors.New("Game not found")
	ErrUserExists       = errors.New("user already exists")
	ErrUserIDRequired   = errors.New("user id required")
	ErrUserNameRequired = errors.New("user's username required")
	ErrMovieIDRequired  = errors.New("movie id required")
	ErrGameIDRequired   = errors.New("Game id required")
	ErrInvalidJSON      = errors.New("invalid json")
	ErrUserRequired     = errors.New("user required")
	ErrGameRequired     = errors.New("game required")
	ErrInvalidEntry     = errors.New("invalid Entry")
)

//login or Signup error
var (
	ErrUserNullPointer  = errors.New("User value is nill or User is Empty")
	ErrUserNotCached    = errors.New("Unable to save User in Cache or Session")
	ErrUserNameEmpty    = errors.New("Username is Empty please enter a Username")
	ErrUsrDbUnreachable = errors.New("Unable to get the UserDB into the Method")
	ErrMovDbUnreachable = errors.New("Unable to get the MovieDB into the Method")
	ErrGamDbUnreachable = errors.New("Unable to get the GameDB into the Method")
)

//Session ErrorTypes
var (
	ErrSessionCookieSaveError = errors.New("could not save cookie session please ensure cookie is enable on your browser")
	ErrIvalidRedirect         = errors.New("invalid redirect URL, Please try again")
	ErrSessionCookieError     = errors.New("could not create a cookie session please ensure cookie is enable on your browser")
)

// Error returns the error message.
//func (e Error) Error() string { return string(e) }
func Error(a ...interface{}) interface{} {
	return a
}
