package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"portfolio/handlers"
	"portfolio/utils"
)

const port = ":80"

func Emailer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		subject := r.FormValue("subject")
		message := r.FormValue("message")

		session, _ := utils.Store.Get(r, "session")

		if email == "" && subject == "" && message == "" {
			session.AddFlash("Try setting a value first.")
		} else {
			auth := smtp.PlainAuth("", "no-reply@coffshire.com", "no-reply password", "smtp.gmail.com")
			to := []string{"jakecoffman@gmail.com"}
			payload := []byte(fmt.Sprintf("Subject: Portfolio contact\r\n\r\nEmail: %s\r\nSubject: %s\r\nMessage:\r\n%s", email, subject, message))
			smtp.SendMail("smtp.gmail.com:587", auth, "no-reply@coffshire.com", to, payload)

			session.AddFlash("Email sent.")
		}
		session.Save(r, w)
	}

	http.Redirect(w, r, "/#contact", http.StatusFound)
}

func Projects(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			http.Error(w, "project not found", http.StatusNotFound)
		}
	}()

	vars := mux.Vars(r)
	project := vars["project"]

	t := template.Must(template.New("base.html").Funcs(utils.FuncMap).ParseFiles(
		"web/templates/base.html",
		"web/templates/projects.html",
		fmt.Sprintf("web/templates/projects/%s.html", project),
		"web/templates/side.html",
	))

	context := make(map[string]string)
	context["project"] = project

	if err := t.Execute(w, context); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	wd, _ := os.Getwd()
	println("Working directory", wd)

	server := &http.Server{Addr: port, Handler: nil}
	l, e := net.Listen("tcp", port)
	if e != nil {
		panic(e)
	}

	r := mux.NewRouter()
	r = mux.NewRouter()

	r.HandleFunc("/", handlers.IndexHandler)
	r.HandleFunc("/contact", Emailer)
	r.HandleFunc("/projects/{project}", Projects)

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	go func() {
		// run stuff outside of the main loop
	}()

	go func() {
		for {
			fmt.Println("Running on " + port)
			server.Serve(l)
			l, e = net.Listen("tcp", port)
			if e != nil {
				panic(e)
			}
		}
	}()
	defer l.Close()

	select {}

	println("Dead")
}
