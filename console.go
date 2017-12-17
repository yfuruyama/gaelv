package gaelv

import (
	"fmt"
	"strings"
	"time"
)

type Console struct {
}

func NewConsole() *Console {
	return &Console{}
}

func (c *Console) PrintLog(r *RequestLog) {
	// format app log
	appLogLines := make([]string, 0, len(r.AppLogs))
	for _, a := range r.AppLogs {
		timestamp := time.Time(a.Time).Format("15:04:05.000")
		line := fmt.Sprintf("    %s %s %s", timestamp, symbolForLevel(a.Level), a.Message)
		appLogLines = append(appLogLines, line)
	}
	appLogStr := strings.Join(appLogLines, "\n")
	if appLogStr != "" {
		appLogStr += "\n" // add trailing newline
	}

	// format entire request log
	level := r.GetLevel()
	timestamp := time.Time(r.StartTime).Format("2006-01-02 15:04:05.000")
	fmt.Printf("%s %s %s %d %s %s %s\n%s", timestamp, symbolForLevel(level), r.Method, r.Status, r.ResponseSizeStr(), r.LatencyStr(), r.Resource, appLogStr)
}

func symbolForLevel(l LogLevel) string {
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
