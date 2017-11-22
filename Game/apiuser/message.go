package apiuser

import (
	"fmt"
	"regexp"
)

type Message struct {
	Email          string
	Password       string
	PasswordVerify string
	Username       string
	//Content      string
	Errors    map[string]interface{}
	Firstname string
	Lastname  string
	ID        string
	LoginURL  string
}

func (msg *Message) Validate() bool {
	reg := regexp.MustCompile(".+@.+\\..+")
	fmt.Println(reg)
	matched := reg.Match([]byte(msg.Email))
	if matched == false {
		msg.Errors["Email"] = "Please enter a valid email"
	}

	if msg.PasswordVerify != msg.Password {
		msg.Errors["Password"] = "Passwords do not match"
	}

	return len(msg.Errors) == 0

}
