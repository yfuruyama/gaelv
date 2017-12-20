package gaelv

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Provider struct {
	lastRequestID int
	db            *sql.DB
}

const (
	POLLING_INTERVAL time.Duration = 100 * time.Millisecond
)

func NewProvider(logsPath string) (*Provider, error) {
	if _, err := os.Stat(logsPath); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("Log file not exist: %s", logsPath))
	}

	db, err := sql.Open("sqlite3", logsPath)
	if err != nil {
		return nil, err
	}

	id, err := FetchLastRequestId(db)
	if err != nil {
		return nil, err
	}

	return &Provider{
		lastRequestID: id,
		db:            db,
	}, nil
}

// block until next log comes in
func (p *Provider) Next() (*RequestLog, error) {
	ticker := time.NewTicker(POLLING_INTERVAL)
	var lastId int
	for {
		<-ticker.C
		id, err := FetchLastRequestId(p.db)
		if err != nil {
			return nil, err
		}
		if id == p.lastRequestID {
			continue
		}
		lastId = id

		ticker.Stop()
		break
	}
	p.lastRequestID = lastId

	return FetchRequestLog(p.db, p.lastRequestID)
}

func (p *Provider) GetLatestLogs(num int) ([]*RequestLog, error) {
	logs := make([]*RequestLog, 0, num)
	for i := num - 1; i >= 0; i-- {
		l, err := FetchRequestLog(p.db, p.lastRequestID-i)
		if err != nil {
			return nil, err
		}
		if l != nil {
			logs = append(logs, l)
		}
	}
	return logs, nil
}
