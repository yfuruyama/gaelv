package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type SSEServer struct {
	logc chan *RequestLog
}

func NewSSEServer(logc chan *RequestLog) *SSEServer {
	return &SSEServer{
		logc: logc,
	}
}

func (s *SSEServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "no-cache")

	for {
		log := <-s.logc
		fmt.Fprintf(w, "data: %s\n\n", log.ToJSON())
		w.(http.Flusher).Flush()
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("error parsing template")
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("error creating html")
	}
}
