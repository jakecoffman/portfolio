package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type Email struct {
	Email   string `json:email`
	Subject string `json:subject`
	Message string `json:message`
}

func Emailer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var data Email
		err := decoder.Decode(&data)
		if err != nil {
			http.Error(w, "Failed loading JSON", http.StatusBadRequest)
			return
		}

		if data.Email == "" && data.Subject == "" && data.Message == "" {
			http.Error(w, "Try setting something first", http.StatusBadRequest)
			return
		} else {
			c, err := smtp.Dial("localhost:25")
			if err != nil {
				println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			c.Mail("no-reply@jakecoffman.com")
			c.Rcpt("jake.coffman@gmail.com")
			wc, err := c.Data()
			if err != nil {
				println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer wc.Close()
			buf := bytes.NewBufferString(fmt.Sprintf(`From: no-reply@jakecoffman.com
To: jake.coffman@gmail.com
Subject: portfolio contact

Email: %s
Subject: %s
Message:
%s`, data.Email, data.Subject, data.Message))
			if _, err = buf.WriteTo(wc); err != nil {
				println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(201)
			return

		}
	}

	http.Redirect(w, r, "/#/?scrollTo=contact", http.StatusFound)
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
