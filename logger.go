package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// Log contructs a LogMessage and dumps the same in chanbuggLog.
// The loglevels are incremental where DEBUG being the highest one and includes all log levels.
// DBGRAM is always logged.
// Arguments:
// component string: name of module / name of webserver
// logLevel LogLevel: DBGRAM, DEBUG, INFO, WARNING, ERROR, FATAL
// message string: format string
// args interface: varargs for message
func Log(component string, logLevel LogLevel, message string, args ...interface{}) {

	programCounter, fileName, lineNumber, _ := runtime.Caller(1)

	if logLevel <= configuredLogLevel || logLevel == DBGRM {
		logMessage := LogMessage{
			TimeStamp:    time.Now(),
			Level:        logLevel.String(),
			Component:    component,
			Message:      message,
			SourceFile:   fileName,
			LineNumber:   lineNumber,
			FunctionName: runtime.FuncForPC(programCounter).Name(),
		}
		chanBuffLog <- logMessage
	}
}

// WebLog contructs a LogMessage and dumps the same in chanBuffLog.
// Arguments:
// component string: name of module / name of webserver
// logLevel LogLevel: DBGRAM, DEBUG, INFO, ERROR, FATAL
// method string: HTTP method e.g. GET, POST, PUT, OPTIONS
// clientIP string: Client IP Address
// path string: URL path
// StatusCode int: HTTP Status Code
// latency time.Time: time required for handling the reques
// message string: log message
func WebLog(component string, logLevel LogLevel, method string, clientIP string, path string, StatusCode int, latency time.Duration, message string) {
	if logLevel <= configuredLogLevel || logLevel == DBGRM {
		logMessage := LogMessage{
			TimeStamp:  time.Now(),
			Level:      logLevel.String(),
			Component:  component,
			StatusCode: StatusCode,
			Latency:    latency,
			ClientIP:   clientIP,
			Method:     method,
			Path:       path,
			Message:    message,
		}
		chanBuffLog <- logMessage
	}
}

// LogDispatcher infinitely waits on channel chanBuffLog,
// extracts data from the channel and dumps log into the file pointed by pServerLogFile.
// Arguments:
// wg *sync.WaitGroup: waitgroup handler for conveying done status to the caller.
// doneChan chan bool: done channel to terminate logger thread.
func LogDispatcher(ploggerWG *sync.WaitGroup, doneChan chan bool) {
	defer func() {
		fmt.Printf("[HEED] Logger Exiting.\n")
		ploggerWG.Done()
	}()

	runFlag := true
	for runFlag {
		select {
		case logMsg, isOK := <-chanBuffLog:
			if !isOK {
				runFlag = false
				break
			}
			dumpServerLog(logMsg)
			break

		case <-doneChan: // chanBuffLog needs to be closed. pull all the logs from the channel and dump them to file-system.
			runFlag = false
			fmt.Printf("[HEED] Flushing log buffer.\n")
			close(chanBuffLog)
			for logMsg := range chanBuffLog {
				dumpServerLog(logMsg)
			}
			break
		}
	}
}

// Initializes logger.
// Arguments:
// logDir string: should be the directory where logs should be generated.
// logLevel LogLevel: describes the severity of log. All logs whose severity is below this level will be discarded.
// Return Values:
// Returns true if logger was successfully initialized, false otherwise.
func Init(logDir string, logLevel LogLevel) bool {
	var err error

	// check if logger is already initialized
	if isInit {
		return true
	}

	// set log level
	configuredLogLevel = logLevel

	// check if log dir exists
	_, err = os.Stat(logDir)
	if err != nil {
		fmt.Printf("Error: Stat(%s): %s", logDir, err)
		return false
	}

	// create buffered channel for logs
	chanBuffLog = make(chan LogMessage, chanbuffLogCapacity)

	// prepare for log rotation
	logfileNameList = make([]string, log_MAX_FILES)
	logFile := filepath.Join(logDir, log_FILE_NAME_PREFIX) + ".1"
	tmpLogFile := filepath.Join(logDir, log_FILE_NAME_PREFIX)
	dummyLogfile = logFile + ".dummy"

	pServerLogFile, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error: OpenFile(%s): %s\n", logFile, err)
		return false
	}

	for i := uint8(0); i < log_MAX_FILES; i++ {
		logfileNameList[i] = fmt.Sprintf("%s.%d", tmpLogFile, i+1)
	}

	errDup2 := syscall.Dup2(int(pServerLogFile.Fd()), syscall.Stdout)
	if errDup2 != nil {
		fmt.Printf("Error: Dup2 - Failed to reuse STDOUT: %s\n", errDup2)
	}

	isInit = true
	return true
}
