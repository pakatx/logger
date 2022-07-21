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

func WebLog(component string, logLevel LogLevel, method string, clientIP string, path string, StatusCode int, latency string, message string) {
	t := time.Now()
	zonename, _ := t.In(time.Local).Zone()
	msgTimeStamp := fmt.Sprintf("%02d-%02d-%d:%02d%02d%02d-%06d-%s",
		t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), zonename)

	// TODO: use t.Format(time.RFC3339Nano) in msgTimeStamp when log aggregator would be used

	if current_LOG_LEVEL == OFF {
		return
	}

	if logLevel <= current_LOG_LEVEL || logLevel == DBGRM {
		logMessage := LogMessage{
			TimeStamp:  msgTimeStamp,
			Level:      logLevel.String(),
			Component:  component,
			StatusCode: StatusCode,
			Latency:    latency,
			ClientIP:   clientIP,
			Method:     method,
			Path:       path,
			Message:    message,
		}
		chanbuffLog <- logMessage
	}
}

// Log contructs a LogMessage and dumps the same in chanbuggLog.
// The loglevels are incremental where DEBUG being the highest one and includes all log levels.
// DBGRAM is always logged.
// Arguments:
// component string: name of module / name of webserver
// logLevel LogLevel: DBGRAM, DEBUG, INFO, ERROR, FATAL, OFF
func Log(component string, logLevel LogLevel, message string, args ...interface{}) {
	t := time.Now()
	zonename, _ := t.In(time.Local).Zone()
	msgTimeStamp := fmt.Sprintf("%02d-%02d-%d:%02d%02d%02d-%06d-%s",
		t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), zonename)

	// TODO: use t.Format(time.RFC3339Nano) in msgTimeStamp when log aggregator would be used

	// TODO: handle 4th return value i.e. OK? Its possible that called may not be able to detect
	// program counter, file and line number as per documentation of runtime.Caller()
	pc, fileName, lineNumber, _ := runtime.Caller(1)

	if current_LOG_LEVEL == OFF {
		return
	}

	if logLevel <= current_LOG_LEVEL || logLevel == DBGRM {
		logMessage := LogMessage{
			TimeStamp:    msgTimeStamp,
			Level:        logLevel.String(),
			Component:    component,
			Message:      message,
			SourceFile:   fileName,
			LineNumber:   lineNumber,
			FunctionName: runtime.FuncForPC(pc).Name(),
		}
		chanbuffLog <- logMessage
	}
}

// LogDispatcher infinitely waits on channel chanbuffLog,
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
		case logMsg, isOK := <-chanbuffLog:
			if !isOK {
				runFlag = false
				break
			}
			dumpServerLog(logMsg)
			break

		case <-doneChan: // chanbuffLog needs to be closed. pull all the logs from the channel and dump them to file-system.
			runFlag = false
			fmt.Printf("[HEED] Flushing log buffer.\n")
			close(chanbuffLog)
			for logMsg := range chanbuffLog {
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
	current_LOG_LEVEL = logLevel // TODO: handle this later in Log and WebLog, if given level is less than current, do not log

	// check if log dir exists
	_, err = os.Stat(logDir)
	if err != nil {
		fmt.Printf("Error: Stat(%s): %s", logDir, err)
		return false
	}

	// create buffered channel for logs
	chanbuffLog = make(chan LogMessage, chanbuffLogCapacity)

	// prepare for log rotation
	logfileNameList = make([]string, log_MAX_FILES)
	logFile := filepath.Join(logDir, log_FILE_NAME_PREFIX) + ".1"
	tmpLogFile := filepath.Join(logDir, log_FILE_NAME_PREFIX)
	dummyLogfile = logFile + ".dummy"

	pServerLogFile, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error: OpenFile(%s): %s\n", logFile, err)
		return false
	}

	for i := uint8(0); i < log_MAX_FILES; i++ {
		logfileNameList[i] = fmt.Sprintf("%s.%d", tmpLogFile, i+1)
	}

	errDup2 := syscall.Dup2(int(pServerLogFile.Fd()), syscall.Stdout) // TODO: check what this does exactly
	if errDup2 != nil {
		fmt.Printf("Error: Dup2 - Failed to reuse STDOUT: %s\n", errDup2)
	}

	isInit = true
	return true
}
