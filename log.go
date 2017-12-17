package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

/*
CREATE TABLE IF NOT EXISTS RequestLogs (
    id INTEGER NOT NULL PRIMARY KEY,
    user_request_id TEXT NOT NULL,
    app_id TEXT NOT NULL,
    version_id TEXT NOT NULL,
    module TEXT NOT NULL,
    ip TEXT NOT NULL,
    nickname TEXT NOT NULL,
    start_time INTEGER NOT NULL,
    end_time INTEGER DEFAULT 0 NOT NULL,
    method TEXT NOT NULL,
    resource TEXT NOT NULL,
    http_version TEXT NOT NULL,
    status INTEGER DEFAULT 0 NOT NULL,
    response_size INTEGER DEFAULT 0 NOT NULL,
    user_agent TEXT NOT NULL,
    url_map_entry TEXT DEFAULT '' NOT NULL,
    host TEXT NOT NULL,
    referrer TEXT,
    task_queue_name TEXT DEFAULT '' NOT NULL,
    task_name TEXT DEFAULT '' NOT NULL,
    latency INTEGER DEFAULT 0 NOT NULL,
    mcycles INTEGER DEFAULT 0 NOT NULL,
    finished INTEGER DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS AppLogs (
    id INTEGER NOT NULL PRIMARY KEY,
    request_id INTEGER NOT NULL,
    timestamp INTEGER NOT NULL,
    level INTEGER NOT NULL,
    message TEXT NOT NULL,
    FOREIGN KEY(request_id) REFERENCES RequestLogs(id)
);
*/

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	CRITICAL
)

type RequestLog struct {
	ID            int
	UserRequestID string
	AppID         string
	VersionID     string
	Module        string
	IP            string
	Nickname      string
	StartTime     LogTime
	EndTime       LogTime
	Method        string
	Resource      string
	HTTPVersion   string
	Status        int32
	ResponseSize  int64
	UserAgent     string
	URLMapEntry   string
	Host          string
	Referrer      string
	TaskQueueName string
	TaskName      string
	Latency       time.Duration
	MCycles       int64
	Finished      bool
	AppLogs       []AppLog
}

type AppLog struct {
	ID      int
	Time    LogTime
	Level   LogLevel
	Message string
}

type LogTime time.Time

func (t *LogTime) Scan(src interface{}) error {
	unixnano := src.(int64)
	sec := int64(unixnano / 1000000)
	nanosec := unixnano - (sec * 1000000)
	*t = LogTime(time.Unix(sec, nanosec))
	return nil
}

func FetchRequestLog(db *sql.DB, id int) (*RequestLog, error) {
	var r RequestLog
	if err := db.QueryRow("SELECT id, start_time, end_time, method, status, response_size, resource FROM RequestLogs WHERE id = ?", id).Scan(&r.ID, &r.StartTime, &r.EndTime, &r.Method, &r.Status, &r.ResponseSize, &r.Resource); err != nil {
		return nil, err
	}

	// var appLogs
	rows, err := db.Query("SELECT id, timestamp, level, message FROM AppLogs WHERE request_id = ?", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var a AppLog
		if err := rows.Scan(&a.ID, &a.Time, &a.Level, &a.Message); err != nil {
			return nil, err
		}
		r.AppLogs = append(r.AppLogs, a)
	}

	return &r, nil
}

func (r *RequestLog) GetLevel() LogLevel {
	if len(r.AppLogs) == 0 {
		return INFO
	}
	level := r.AppLogs[0].Level
	for _, a := range r.AppLogs {
		if a.Level > level {
			level = a.Level
		}
	}
	return level
}

func (r *RequestLog) LatencyStr() string {
	latencyNanos := time.Time(r.EndTime).Sub(time.Time(r.StartTime)).Nanoseconds()
	// over 1 sec.
	if latencyNanos/1000000000 > 0 {
		return fmt.Sprintf("%0.1fs", float32(latencyNanos/1000000000))
	} else {
		return fmt.Sprintf("%dms", latencyNanos/1000000)
	}
}

func (r *RequestLog) ResponseSizeStr() string {
	// over 1 KB
	if r.ResponseSize >= 1024 {
		return fmt.Sprintf("%0.1fKB", float32(r.ResponseSize/1024))
	} else {
		return fmt.Sprintf("%dB", r.ResponseSize)
	}
}

func (r *RequestLog) ToJSON() string {
	j, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	return string(j)
}
