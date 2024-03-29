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
var chanBuffLog chan LogMessage

// global flag to restrict reinitiation of logger
var isInit bool

// configured log level; set to DEBUG by default
var configuredLogLevel LogLevel = DEBUG
