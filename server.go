package gaelv

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type SSEServer struct {
	connected    chan chan *RequestLog
	disconnected chan chan *RequestLog
	clients      map[chan *RequestLog]bool
	Logc         chan *RequestLog
}

func NewSSEServer() *SSEServer {
	return &SSEServer{
		connected:    make(chan (chan *RequestLog)),
		disconnected: make(chan (chan *RequestLog)),
		clients:      make(map[chan *RequestLog]bool),
		Logc:         make(chan *RequestLog),
	}
}

func (s *SSEServer) Start() {
	go func() {
		for {
			select {
			case client := <-s.connected:
				s.clients[client] = true
				log.Println("client connected")
			case client := <-s.disconnected:
				delete(s.clients, client)
				close(client)
				log.Println("client disconnected")
			case l := <-s.Logc:
				if len(s.clients) > 0 {
					for client, _ := range s.clients {
						client <- l
					}
					log.Println("sent new log")
				}
			}
		}
	}()
}

func (s *SSEServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "no-cache")

	c := make(chan *RequestLog)
	s.connected <- c

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		s.disconnected <- c
	}()

	for {
		l, ok := <-c
		if !ok {
			// client disconnected
			break
		}
		fmt.Fprintf(w, "data: %s\n\n", l.ToJSON())
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
