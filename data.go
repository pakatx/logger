package logger

import (
	"os"
)

// active log file handler
var pServerLogFile *os.File

// log rotation
var currentLogfileCnt uint8 = 1
var logfileNameList []string
var dummyLogfile string

// buffered channel for logs
var chanbuffLog chan LogMessage

// default log level setting
var current_LOG_LEVEL LogLevel = DBGRM

// global flag to restrict reinitiation of logger
var isInit bool
