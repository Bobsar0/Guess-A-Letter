package apiuser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Bobsar0/Guess-A-Letter/tools"

	"github.com/go-chi/chi"
)

// UserHandler represents an HTTP API handler for users.
type UserHandler struct {
	mux *chi.Mux
	//redirectURL string
	Session *Session
	//user        *User
	Logger *log.Logger
}

// NewUserHandler returns a new instance of UserHandler.
func NewUserHandler(s *Session) *UserHandler {
	h := &UserHandler{
		mux:     chi.NewRouter(),
		Logger:  log.New(os.Stderr, "", log.LstdFlags),
		Session: s,
	}
	h.mux.Get("/", s.Validate(h.indexHandler))
	h.mux.Get("/signup", h.signupHandler)
	h.mux.Post("/signup", h.signupPostHandler)
	h.mux.Get("/login", s.Validate(h.loginHandler))
	h.mux.Post("/login", h.loginPostHandler)
	h.mux.Post("/logout", h.logoutHandler)
	h.mux.Get("/users/list", s.Validate(h.handleListUsers))
	h.mux.Get("/users/get/:username", s.Validate(h.handleGetUser))
	h.mux.Post("/users/update/:username", s.Validate(h.handleUpdateUser))
	return h
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *UserHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		user = nil
	} else if user.ImageURL == "" {
		user.ImageURL = "/tools/asset/images/user.png"
	}
	if err := tools.IndexTmpl.Execute(w, r, nil, user, true); err != nil {
		fmt.Errorf("%v, %v, %v, %v", w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	return
}

func (h *UserHandler) signupHandler(w http.ResponseWriter, r *http.Request) {

	msg := &Message{
		LoginURL: "/login?redirect=/", // + r.URL.RequestURI(),
	}
	tools.SignupTmpl.Execute(w, r, msg, nil, true)
	return
}

func (h *UserHandler) signupPostHandler(w http.ResponseWriter, r *http.Request) {
	msg := &Message{
		Email:          r.FormValue("email"),
		Firstname:      r.FormValue("firstname"),
		Lastname:       r.FormValue("lastname"),
		Password:       r.FormValue("password"),
		PasswordVerify: r.FormValue("passwordVerify"),
		Username:       r.FormValue("username"),
	}
	msg.Errors = make(map[string]interface{})
	// username taken or password or email incorrect ?
	_, err := h.Session.Userservicesess.GetUser(msg.Username) //get user with username and check if it exist
	if err == nil || err == tools.ErrUserNameEmpty || msg.Validate() == false {
		if err == nil {
			msg.Errors["Username"] = fmt.Sprintf("Username \"%s\" is already taken, please try another", msg.Username)
			_ = msg.Validate()
		} else if err == tools.ErrUserNameEmpty {
			msg.Errors["Username"] = tools.ErrUserNameEmpty
			_ = msg.Validate()
		}
		if err := tools.SignupTmpl.Execute(w, r, msg, nil, true); err != nil {
			tools.Error(w, err, http.StatusInternalServerError, h.Logger)
			return
		}
		return
	}
	// get new valid values form request values
	usr := &User{
		Email:     r.FormValue("email"),
		Username:  r.FormValue("username"),
		Password:  r.FormValue("password"),
		Firstname: r.FormValue("firstname"),
		Lastname:  r.FormValue("lastname"),
	}
	usr.DisplayName = usr.Firstname + " " + usr.Lastname
	if err := h.Session.Userservicesess.AddUser(usr); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	//Set token for the user in cookie
	h.Session.setToken(w, r, usr.Username, "/", usr.Level)

	return
}

func (h *UserHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		if err := tools.LoginTmpl.Execute(w, r, nil, nil, true); err != nil {
			tools.Error(w, err, http.StatusInternalServerError, h.Logger)
			return
		}
		return
	}
	if err := tools.IndexTmpl.Execute(w, r, nil, user, true); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	return
}

func (h *UserHandler) loginPostHandler(w http.ResponseWriter, r *http.Request) {

	msg := &Message{
		Password: r.FormValue("password"),
		//PasswordVerify: r.FormValue("passwordVerify"),
		Username: r.FormValue("username"),
	}
	msg.Errors = make(map[string]interface{})
	// username taken or password or email incorrect ?
	user, err := h.Session.Userservicesess.GetUser(msg.Username) //get user with username and check if it exist
	if err != nil {
		if err == tools.ErrUnauthorized {
			msg.Errors["Username"] = fmt.Sprintf("Unknown Username or Password, Try and signup")
		} else if err == tools.ErrUserNameEmpty {
			msg.Errors["Username"] = tools.ErrUserNameEmpty
		} else {
			msg.Errors["Username"] = fmt.Sprintf("Unknown Username or Password, Try and signup")
		}
		if err := tools.LoginTmpl.Execute(w, r, msg, nil, true); err != nil {
			tools.Error(w, err, http.StatusInternalServerError, h.Logger)
			return
		}
		return
	}
	//Set token for the user in cookie
	h.Session.setToken(w, r, user.Username, "/", user.Level)
	return
}

// deletes the cookie
func (h *UserHandler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	h.Session.logout(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func (h *UserHandler) handleListUsers(w http.ResponseWriter, r *http.Request) {
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		//ToDo :prompt user of enabling cookies
		tools.Error(w, errors.New("Your are not an Admin"), http.StatusInternalServerError, h.Logger)
		return
	}
	if user.Level == "Admin" {
		if user.ImageURL == "" {
			user.ImageURL = "/tools/asset/images/user.png"
		}
		usrs, err := h.Session.Userservicesess.ListUsers()
		if err != nil {
			tools.Error(w, err, http.StatusInternalServerError, h.Logger)
			return
		}
		dat := struct {
			Usrs       []*User
			AddUserURL string
		}{
			Usrs:       usrs,
			AddUserURL: "/Users/add?redirect=/", //the resaon is for list page to know if user is loged in
		}
		if err := tools.ListusersTmpl.Execute(w, r, dat, user, true); err != nil {
			tools.Error(w, err, http.StatusInternalServerError, h.Logger)
			return
		}
	}
	//tools.Error(w, errors.New("Your are not an Admin"), http.StatusInternalServerError, h.Logger)
	return
}

type addUserRequest struct {
	User  *User  `json:"user,omitempty"`
	Token string `json:"token,omitempty"`
}
type addUserResponse struct {
	User *User  `json:"user,omitempty"`
	Err  string `json:"err,omitempty"`
}

// handleGetUser handles requests to fetch a single user.
func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	// Find user by ID.
	d, err := h.Session.Userservicesess.GetUser(username)
	if err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
	} else if d == nil {
		tools.Error("Notfound", w)
	} else {
		encodeJSON(w, &getUserResponse{User: d}, h.Logger)
	}
}

type getUserResponse struct {
	User *User  `json:"user,omitempty"`
	Err  string `json:"err,omitempty"`
}

// handleupdateUser handles requests to update a user level.
func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Decode request.
	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		tools.Error(w, tools.ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	username := chi.URLParam(r, "username")

	d := req.User
	d.Token = req.Token
	d.Username = username
	// Create user.
	switch err := h.Session.Userservicesess.UpdateUser(req.User); err {
	case nil:
		encodeJSON(w, &updateUserResponse{}, h.Logger)
	case tools.ErrUserNotFound:
		tools.Error(w, err, http.StatusNotFound, h.Logger)
	case tools.ErrUnauthorized, tools.ErrUserNameRequired:
		tools.Error(w, err, http.StatusUnauthorized, h.Logger)
	default:
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

type updateUserRequest struct {
	User  *User  `json:"user,omitempty"`
	Token string `json:"token"`
}

type updateUserResponse struct {
	Err string `json:"err,omitempty"`
}

// encodeJSON encodes v to w in JSON format. Error() is called if encoding fails.
func encodeJSON(w http.ResponseWriter, v interface{}, logger *log.Logger) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, logger)
	}
}
func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n Inside GetUsersHandler")
	data := struct {
		User  *User  `json:"user,omitempty"`
		Token string `json:"token"`
	}{
		User: &User{
			Username: "Steve",
			Password: "123",
		},
		Token: "sdfws",
	}
	webWriteJSON(w, r, data)
}
func webWriteJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	var u *url.URL
	u = &url.URL{}
	u.Path = "/api/users"
	reqBody, err := json.Marshal(data)
	fmt.Println("\n Iside webWriteJSON", data)
	resp, err := http.Post(u.String(), "application/json", bytes.NewReader(reqBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

// ValidateRedirectURL checks that the URL provided is valid.
// If the URL is missing, redirect the user to the application's root.
// The URL must not be absolute (i.e., the URL must refer to a path within this
// application).
func ValidateRedirectURL(path string) (string, error) {
	if path == "" {
		return "/", nil
	}

	// Ensure redirect URL is valid and not pointing to a different server.
	parsedURL, err := url.Parse(path)
	if err != nil {
		return "/", err
	}
	if parsedURL.IsAbs() {
		return "/", errors.New("URL invalid: URL must not be absolute")
	}
	if strings.Contains(path, "/signup?redirect=") {
		return path[17:], nil
	}
	return path, nil
}
