package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"log"

	"github.com/addsict/gaelv"
)

var usage = `Usage: gaelv [options...]

Options:
    --logs_path        Path to logs
    --no-server        Don't serve http server
`

func main() {
	var logsPath string
	var port int

	flag.StringVar(&logsPath, "logs_path", "", "logs path")
	flag.IntVar(&port, "port", 9090, "server port")
	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	flag.Parse()

	provider := gaelv.NewProvider(logsPath)
	logc := make(chan *gaelv.RequestLog)
	console := gaelv.NewConsole()
	go func() {
		for {
			requestLog, err := provider.Next()
			if err != nil {
				log.Fatal(err)
			}
			logc <- requestLog // TODO: blocked?
			console.PrintLog(requestLog)
		}
	}()

	s := gaelv.NewSSEServer(logc)
	http.Handle("/event/logs", s)
	http.Handle("/", http.HandlerFunc(gaelv.IndexHandler))

	log.Printf("Starting log viewer at: http://localhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
