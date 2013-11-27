package main

import (
	"fmt"
	"github.com/jakecoffman/portfolio/utils"
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
				w.Write([]byte(err.Error()))
				return
			} else {
				http.Redirect(w, r, "/email-sent.html", http.StatusFound)
				return
			}
		}
		session.Save(r, w)
	}

	http.Redirect(w, r, "/#contact", http.StatusFound)
}

func main() {
	if password == "" {
		println("*** 'emailpassword' env var not set, emails won't be sent ***")
	}

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

	http.HandleFunc("/contact", Emailer)
	http.Handle("/", http.FileServer(http.Dir("web/")))

	server := &http.Server{Addr: port, Handler: nil}
	l, e := net.Listen("tcp", port)
	if e != nil {
		panic(e)
	}
	defer l.Close()

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
