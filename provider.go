package gaelv

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Provider struct {
	lastRequestID int
	db            *sql.DB
}

const (
	POLLING_INTERVAL time.Duration = 500 * time.Millisecond
)

func NewProvider(logsPath string) *Provider {
	db, err := sql.Open("sqlite3", logsPath)
	if err != nil {
		log.Fatal(err)
	}

	return &Provider{
		lastRequestID: 0,
		db:            db,
	}
}

// block until next log comes in
func (p *Provider) Next() (*RequestLog, error) {
	ticker := time.NewTicker(POLLING_INTERVAL)
	var id int
	for {
		<-ticker.C
		// finished=0 means the request is still processing
		if err := p.db.QueryRow("SELECT id FROM RequestLogs WHERE id > ? AND finished = 1", p.lastRequestID).Scan(&id); err != nil {
			switch {
			case err == sql.ErrNoRows:
				continue
			default:
				return nil, err
			}
		}
		ticker.Stop()
		break
	}
	// log.Printf("New log found: %d\n", id)
	p.lastRequestID = id

	return FetchRequestLog(p.db, id)
}
