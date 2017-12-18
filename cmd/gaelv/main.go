package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"log"

	"github.com/addsict/gaelv"
)

var usage = `Usage:

    gaelv --logs_path=</path/to/log.db> [options...]

Options:

    --logs_path        Path to logs file
    --port             Port for server
    --no_server        Stop running http server (console mode)
`

func main() {
	var logsPath string
	var port int
	var noServer bool

	flag.StringVar(&logsPath, "logs_path", "", "Path to logs file")
	flag.IntVar(&port, "port", 9090, "Port for server")
	flag.BoolVar(&noServer, "no_server", false, "Stop running server (console mode)")
	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	flag.Parse()

	if logsPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	provider, err := gaelv.NewProvider(logsPath)
	if err != nil {
		log.Fatal(err)
	}

	if noServer {
		console := gaelv.NewConsole()
		for {
			requestLog, err := provider.Next()
			if err != nil {
				log.Fatal(err)
			}
			console.PrintLog(requestLog)
		}
	} else {
		logc := make(chan *gaelv.RequestLog)
		go func() {
			for {
				requestLog, err := provider.Next()
				if err != nil {
					log.Fatal(err)
				}
				logc <- requestLog
			}
		}()

		s := gaelv.NewSSEServer(logc)
		http.Handle("/event/logs", s)
		http.Handle("/", http.HandlerFunc(gaelv.IndexHandler))

		log.Printf("Starting log viewer at: http://localhost:%d\n", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}
}
