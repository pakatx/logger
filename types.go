package logger

import (
	"fmt"
	"time"
)

type LogLevel int

func (logLevel LogLevel) String() string {
	switch logLevel {
	case OFF:
		return "OFF"
	case FATAL:
		return "FATAL"
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARNING"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case DBGRM:
		return "DBGRM"
	default:
		return "invalid"
	}
}

type LogMessage struct {
	// common
	Level     string    `json:"level,omitempty"`
	TimeStamp time.Time `json:"time_stamp,omitempty"`
	Component string    `json:"component,omitempty"`
	// context based
	SourceFile   string `json:"source_file,omitempty"`
	LineNumber   int    `json:"line_number,omitempty"`
	FunctionName string `json:"function_name,omitempty"`
	// webserver
	Method     string        `json:"method,omitempty"`
	Path       string        `json:"path,omitempty"`
	ClientIP   string        `json:"client_ip,omitempty"`
	StatusCode int           `json:"status_code,omitempty"`
	Latency    time.Duration `json:"latency,omitempty"`
	// common
	Message string `json:"message,omitempty"`
}

func (logMessage LogMessage) String() string {
	t := logMessage.TimeStamp
	zonename, _ := t.In(time.Local).Zone()
	timeStamp := fmt.Sprintf("%02d-%02d-%d:%02d%02d%02d-%06d-%s", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), zonename)

	latency := logMessage.Latency.Microseconds()

	output := fmt.Sprintf("[%s] [%s] [%s] ", logMessage.Component, timeStamp, logMessage.Level)
	if logMessage.SourceFile != "" {
		output = output + fmt.Sprintf("[%s: %d] [%s]:\n", logMessage.SourceFile, logMessage.LineNumber, logMessage.FunctionName)
	}
	if logMessage.Method != "" {
		output = output + fmt.Sprintf("[%s] [%d] [%s] [%s] [%d]:\n", logMessage.Method, logMessage.StatusCode, logMessage.Path, logMessage.ClientIP, latency)
	}
	output = output + fmt.Sprintf("%s\n", logMessage.Message)
	return output
}
