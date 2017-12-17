package main

import (
	"database/sql"
	"fmt"
	"strings"
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

// 2017-12-17 19:33:47.017 JST POST 204 0 B 1.1 s CloudPubSub-Google /_ah/push-handlers/
func (r *RequestLog) Format() string {
	// format app log
	appLogLines := make([]string, 0, len(r.AppLogs))
	for _, a := range r.AppLogs {
		timestamp := time.Time(a.Time).Format("15:04:05.000")
		line := fmt.Sprintf("    %s %s %s", timestamp, a.Level.symbol(), a.Message)
		appLogLines = append(appLogLines, line)
	}
	appLogStr := strings.Join(appLogLines, "\n")
	if appLogStr != "" {
		appLogStr += "\n" // add trailing newline
	}

	// format entire request log
	level := r.GetLevel()
	timestamp := time.Time(r.StartTime).Format("2006-01-02 15:04:05.000")
	return fmt.Sprintf("%s %s %s %d %s %s %s\n%s", timestamp, level.symbol(), r.Method, r.Status, r.ResponseSizeStr(), r.LatencyStr(), r.Resource, appLogStr)
}

func (l LogLevel) symbol() string {
	switch l {
	case DEBUG:
		return "[D]"
	case INFO:
		return withCyan("[I]")
	case WARNING:
		return withYellow("[W]")
	case ERROR:
		return withRed("[E]")
	case CRITICAL:
		return withMagenta("[C]")
	default:
		return ""
	}
}

func withCyan(s string) string {
	return fmt.Sprintf("\x1b[36m%s\x1b[0m", s)
}

func withYellow(s string) string {
	return fmt.Sprintf("\x1b[33m%s\x1b[0m", s)
}

func withRed(s string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", s)
}

func withMagenta(s string) string {
	return fmt.Sprintf("\x1b[35m%s\x1b[0m", s)
}
