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
	provider     *Provider
}

func NewSSEServer(provider *Provider) *SSEServer {
	return &SSEServer{
		connected:    make(chan (chan *RequestLog)),
		disconnected: make(chan (chan *RequestLog)),
		clients:      make(map[chan *RequestLog]bool),
		Logc:         make(chan *RequestLog),
		provider:     provider,
	}
}

func (s *SSEServer) Start() {
	go func() {
		for {
			select {
			case client := <-s.connected:
				s.clients[client] = true
				log.Println("client connected")

				// send latest logs to the connected client
				logs, err := s.provider.GetLatestLogs(100)
				if err != nil {
					log.Fatal(err)
				}
				for _, l := range logs {
					client <- l
				}
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
	if strings.HasPrefix(r.URL.Path, "/static/") {
		data, err := Asset(strings.TrimLeft(r.URL.Path, "/"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if strings.HasSuffix(r.URL.Path, "css") {
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		}
		if strings.HasSuffix(r.URL.Path, "js") {
			w.Header().Set("Content-Type", "application/javascript")
		}
		fmt.Fprintln(w, string(data))
		return
	}

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := Asset("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.New("").Parse(string(data))
	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Fatalf("error creating html: %s", err)
	}
}
