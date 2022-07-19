package logger

import (
	"os"
)

// log levels
const (
	DBGRM LogLevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
)

const chanbuffLogCapacity int = 10

// log-file file handler.
var pServerLogFile *os.File

var currentLogfileCnt uint8 = 1
var logfileNameList []string
var dummyLogfile string

//var loggerWG sync.WaitGroup

const log_MAX_FILES int8 = 10
const log_FILE_NAME_PREFIX string = "server.log"
const log_FILE_SIZE int64 = 20971520 // 20 MB

var srcBaseDir string
