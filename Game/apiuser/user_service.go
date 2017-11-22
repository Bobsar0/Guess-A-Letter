//1st code
package apiuser

import (
	"Game/tools"
	"errors"
	"fmt"
	"time"
)

type DBType map[string]*User

type UserServiceSess struct {
	session *Session
}

var _ UserService = &UserServiceSess{}

func (u *UserServiceSess) AddUser(user *User) error {

	user.CreatedDate = time.Now()
	db, ok := u.session.db.(DBType)
	if !ok {
		return tools.ErrUsrDbUnreachable
	}
	db[user.Username] = user
	return nil
}
func (u *UserServiceSess) GetUser(username string) (*User, error) {
	db, ok := u.session.db.(DBType)
	if !ok {
		return nil, tools.ErrUsrDbUnreachable
	}
	if username != "" {
		user := db[username]
		if user == nil {
			return nil, tools.ErrUserNotFound
		}
		return user, nil
	}
	return nil, tools.ErrUserNameEmpty
}

func (u *UserServiceSess) DeleteUser(username string) error {

	db, ok := u.session.db.(DBType)
	if !ok {
		return tools.ErrUsrDbUnreachable
	}
	user, ok := db[username]
	if !ok {
		return tools.ErrUserNotFound
	} else if user.Username != username {
		return tools.ErrUnauthorized
	}
	delete(db, username)
	return nil
}

func (u *UserServiceSess) ListUsers() ([]*User, error) {
	db, ok := u.session.db.(DBType)
	if !ok {
		return nil, tools.ErrUsrDbUnreachable
	}
	var usersList []*User
	for _, user := range db {
		usersList = append(usersList, user)
	}
	return usersList, nil
}

func (u *UserServiceSess) UpdateUser(user *User) error {

	db, ok := u.session.db.(DBType)
	if !ok {
		return tools.ErrUsrDbUnreachable
	}
	// Only allow owner to update Product.
	userInDB, ok := db[user.Username]
	if !ok {
		return fmt.Errorf("memory db: product not found with ID %v", user.UserID)
	} else if userInDB.Username != string(user.Username) {
		return tools.ErrUnauthorized
	}
	if user.Username == "" {
		return errors.New("memory db: product with unassigned ID passed into updateProduct")
	}
	db[user.Username] = user
	return nil
}
