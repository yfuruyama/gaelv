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
    --console          Print logs in the console
`

func main() {
	var logsPath string
	var port int
	var consoleMode bool

	flag.StringVar(&logsPath, "logs_path", "", "Path to logs file")
	flag.IntVar(&port, "port", 9090, "Port for server")
	flag.BoolVar(&consoleMode, "console", false, "Print logs in the console")
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

	if consoleMode {
		console := gaelv.NewConsole()
		for {
			requestLog, err := provider.Next()
			if err != nil {
				log.Fatal(err)
			}
			console.PrintLog(requestLog)
		}
	} else {
		s := gaelv.NewSSEServer(provider)
		s.Start()

		go func() {
			for {
				requestLog, err := provider.Next()
				if err != nil {
					log.Fatal(err)
				}
				s.Logc <- requestLog
			}
		}()

		http.Handle("/event/logs", s)
		http.Handle("/", http.HandlerFunc(gaelv.IndexHandler))

		log.Printf("Starting log viewer at: http://localhost:%d\n", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}
}
