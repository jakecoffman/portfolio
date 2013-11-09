package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jakecoffman/portfolio/handlers"
	"github.com/jakecoffman/portfolio/utils"
	"html/template"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"os/signal"
	"time"
)

const port = ":80"

var password = os.Getenv("emailpassword")

var logFile *os.File

var router *mux.Router

func Emailer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		subject := r.FormValue("subject")
		message := r.FormValue("message")

		session, _ := utils.Store.Get(r, "session")

		if email == "" && subject == "" && message == "" {
			session.AddFlash("Try setting a value first.")
		} else {
			auth := smtp.PlainAuth("", "no-reply@coffshire.com", password, "smtp.gmail.com")
			to := []string{"jakecoffman@gmail.com"}
			payload := []byte(fmt.Sprintf("Subject: Portfolio contact\r\n\r\nEmail: %s\r\nSubject: %s\r\nMessage:\r\n%s", email, subject, message))
			err := smtp.SendMail("smtp.gmail.com:587", auth, "no-reply@coffshire.com", to, payload)
			if err != nil {
				session.AddFlash("Failed to send email")
			} else {
				session.AddFlash("Email sent.")
			}
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

func handler(w http.ResponseWriter, r *http.Request) {
	// Before
	router.ServeHTTP(w, r)
	// After
	fmt.Fprintf(logFile, "%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
}

func main() {
	wd, _ := os.Getwd()
	println("Working directory", wd)

	var err error
	logFile, err = os.Create("logs.txt")
	if err != nil {
		log.Fatal("Tried to create log file: ", err)
		return
	}
	defer logFile.Close()

	rand.Seed(time.Now().UTC().UnixNano())

	server := &http.Server{Addr: port, Handler: nil}
	l, e := net.Listen("tcp", port)
	if e != nil {
		panic(e)
	}
	defer l.Close()

	router = mux.NewRouter()
	router.HandleFunc("/", handlers.IndexHandler)
	router.HandleFunc("/contact", Emailer)
	router.HandleFunc("/projects/{project}", Projects)
	router.PathPrefix("/static/").Handler(http.FileServer(http.Dir("web/")))

	http.HandleFunc("/", handler)

	go func() {
		fmt.Println("Running on " + port)
		server.Serve(l)
	}()

	// Capture keyboard interrupt and then stop gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		fmt.Println("Received interrupt, closing listener")
		l.Close()
	}

	println("Dead")
}
