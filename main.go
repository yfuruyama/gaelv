package main

import (
	"fmt"
	"net/http"

	"log"
)

func main() {
	provider := NewProvider("/tmp/gaelog.db")
	logc := make(chan *RequestLog)
	console := NewConsole()
	go func() {
		for {
			requestLog, err := provider.Next()
			if err != nil {
				log.Fatal(err)
			}
			go func() {
				logc <- requestLog
			}()
			console.PrintLog(requestLog)
		}
	}()

	s := NewSSEServer(logc)
	http.Handle("/event/logs", s)
	http.Handle("/", http.HandlerFunc(IndexHandler))

	port := 9000
	log.Printf("Starting log viewer at: http://localhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
