package gaelv

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
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
	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	if strings.HasPrefix(r.URL.Path, "/static/") {
		fileServer.ServeHTTP(w, r)
		return
	}

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Fatalf("error creating html: %s", err)
	}
}
