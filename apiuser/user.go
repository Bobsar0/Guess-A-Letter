package apiuser

import "time"

type User struct {
	Username      string
	Password      string
	Firstname     string
	Lastname      string
	UserID        string
	Email         string
	DisplayName   string
	ImageURL      string
	Token         string
	Url           string
	Authenticated bool
	CreatedDate   time.Time
	Expiry        int64
	Level         string
}
type UserService interface {
	AddUser(*User) error
	GetUser(username string) (*User, error)
	DeleteUser(username string) error
	ListUsers() ([]*User, error)
	UpdateUser(*User) error
}
