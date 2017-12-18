package gaelv

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
	ID            int           `json:"id"`
	UserRequestID string        `json:"userRequestId"`
	AppID         string        `json:"appId"`
	VersionID     string        `json:"versionId"`
	Module        string        `json:"module"`
	IP            string        `json:"ip"`
	Nickname      string        `json:"nickname"`
	StartTime     LogTime       `json:"startTime"`
	EndTime       LogTime       `json:"endTime"`
	Method        string        `json:"method"`
	Resource      string        `json:"resource"`
	HTTPVersion   string        `json:"httpVersion"`
	Status        int32         `json:"status"`
	ResponseSize  int64         `json:"responseSize"`
	UserAgent     string        `json:"userAgent"`
	URLMapEntry   string        `json:"urlMapEntry"`
	Host          string        `json:"host"`
	Referrer      string        `json:"referrer"`
	TaskQueueName string        `json:"taskQueueName"`
	TaskName      string        `json:"taskName"`
	Latency       time.Duration `json:"latency"`
	MCycles       int64         `json:"mCycles"`
	Finished      bool          `json:"finished"`
	AppLogs       []AppLog      `json:"appLogs"`
}

type AppLog struct {
	ID      int      `json:"id"`
	Time    LogTime  `json:"time"`
	Level   LogLevel `json:"level"`
	Message string   `json:"message"`
}

type LogTime time.Time

func (t *LogTime) Scan(src interface{}) error {
	usec := src.(int64)
	sec := int64(usec / 1e6)
	nanosec := (usec * 1e3) - (sec * 1e9)
	*t = LogTime(time.Unix(sec, nanosec))
	return nil
}

func (t LogTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%f", float64(time.Time(t).UnixNano())/1e9)), nil
}

func FetchRequestLog(db *sql.DB, id int) (*RequestLog, error) {
	var r RequestLog
	if err := db.QueryRow("SELECT id, start_time, end_time, method, status, response_size, resource FROM RequestLogs WHERE id = ?", id).Scan(&r.ID, &r.StartTime, &r.EndTime, &r.Method, &r.Status, &r.ResponseSize, &r.Resource); err != nil {
		return nil, err
	}

	// As of 2017-12-17 latency column in the RequestLogs table is not updated
	r.Latency = time.Duration(time.Time(r.EndTime).Sub(time.Time(r.StartTime)).Nanoseconds())

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
	// over 1 sec.
	if r.Latency/1e9 > 0 {
		return fmt.Sprintf("%0.1f s", float64(r.Latency/1e9))
	} else {
		return fmt.Sprintf("%d ms", r.Latency/1e6)
	}
}

func (r *RequestLog) ResponseSizeStr() string {
	// over 1 KB
	if r.ResponseSize >= 1024 {
		return fmt.Sprintf("%0.1f KB", float64(r.ResponseSize/1024))
	} else {
		return fmt.Sprintf("%d B", r.ResponseSize)
	}
}

func (r *RequestLog) ToJSON() string {
	j, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	return string(j)
}

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case CRITICAL:
		return "CRITICAL"
	default:
		return ""
	}
}

func (l LogLevel) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}
