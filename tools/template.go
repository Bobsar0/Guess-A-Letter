package tools

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

var (
	// See template.go
	SignupTmpl    = parseTemplate("signup.html")
	IndexTmpl     = parseTemplate("index.html")
	LoginTmpl     = parseTemplate("login.html")
	ListusersTmpl = parseTemplate("listusers.html")
	ListgamesTmpl = parseTemplate("listgames.html")
	GameTmpl      = parseTemplate("game.html")
)

// parseTemplate applies a given file to the body of the base template.
func parseTemplate(filename string) *appTemplate {
	tmpl := template.Must(template.ParseFiles("tools/templates/base.html"))

	// Put the named file into a template called "body"
	path := filepath.Join("tools/templates", filename)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("could not read template: %v", err))
	}
	template.Must(tmpl.New("body").Parse(string(b)))

	return &appTemplate{tmpl.Lookup("base.html")}
}

// appTemplate is a user login-aware wrapper for a html/template.
type appTemplate struct {
	t *template.Template
}

// Execute writes the template using the provided data, adding login and user
// information to the base template.
func (tmpl *appTemplate) Execute(w http.ResponseWriter, r *http.Request, dat interface{}, usr interface{}, noFooter bool) error {
	d := struct {
		Data        interface{}
		AuthEnabled bool
		LoginURL    string
		LogoutURL   string
		AddFooter   bool
		SignupURL   string
		User        interface{}
	}{
		Data:        dat,
		AuthEnabled: true,
		LoginURL:    "/login?redirect=" + r.URL.RequestURI(),
		LogoutURL:   "/logout?redirect=" + r.URL.RequestURI(),
		SignupURL:   "/signup?redirect=" + r.URL.RequestURI(),
		AddFooter:   noFooter,
		User:        usr,
	}

	if err := tmpl.t.Execute(w, d); err != nil {
		log.Fatalln(err, "could not write template: %+v", err)
	}
	return nil
}
