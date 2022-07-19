package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

var chanbuffLog chan logmessage

var current_LOG_LEVEL LogLevel = DBGRM

var isInit bool

/* ****************************************************************************
Description :
- Constructs a type logmessage variable.
- Dumps the same in the logmsg_buffered_channel

Arguments   :
1> strcomponent string: Modulename.
2> loglevelStr string:
- There exist 4 loglevels: ERROR, WARNING, INFO, and DEBUG.
The loglevels are incremental where DEBUG being the highest one and
includes all log levels.

Return value: na

Additional note: na
**************************************************************************** */
func Log(strcomponent string, loglevelStr string, msg string, args ...interface{}) {
	t := time.Now()
	zonename, _ := t.In(time.Local).Zone()
	msgTimeStamp := fmt.Sprintf("%02d-%02d-%d:%02d%02d%02d-%06d-%s",
		t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), zonename)

	pc, fn, line, _ := runtime.Caller(1)

	filePath := strings.Split(fn, srcBaseDir)
	srcFile := srcBaseDir + filePath[len(filePath)-1]

	msgPrefix := ""
	if loglevelStr == "DBGRM" {
		msgPrefix = "#### "
	}

	logMsg := fmt.Sprintf("[%s] [%s] [%s] [%s: %d] [%s]:\n", strcomponent, msgTimeStamp, loglevelStr, srcFile, line, runtime.FuncForPC(pc).Name())
	logMsg = fmt.Sprintf(logMsg+msg, args...)
	logMsg = msgPrefix + logMsg + "\n"

	logMessage := logmessage{
		component: strcomponent,
		logmsg:    logMsg,
	}

	chanbuffLog <- logMessage
}

/* ****************************************************************************
Description :
- A go routine, invoked through Logger()
- Waits onto buffered channel name chanbuffLog infinitely.
- Extracts data from the channel, it's of type logmessage.
- Dumps log into the file pointed by pServerLogFile.

Arguments   : na for now.
1> wg *sync.WaitGroup: waitgroup handler for conveying done status to the caller.
2> doneChan chan bool: done channel to terminate logger thread.

Return Value: na

Additional note: na
**************************************************************************** */
func LogDispatcher(ploggerWG *sync.WaitGroup, doneChan chan bool) {
	defer func() {
		fmt.Println("logger exiting.")
		ploggerWG.Done()
	}()

	/* for {
	    select {
	        case logMsg := <-chanbuffLog: // pushes dummy logmessage onto the channel
	            dumpServerLog(logMsg.logmsg)
	    }
	} */

	runFlag := true
	for runFlag {
		select {
		case logMsg, isOK := <-chanbuffLog: // pushes dummy logmessage onto the channel
			if !isOK {
				runFlag = false
				break
			}
			dumpServerLog(logMsg.logmsg)
			break

		case <-doneChan: // chanbuffLog needs to be closed. pull all the logs from the channel and dump them to file-system.
			runFlag = false
			dumpServerLog("[WARNING]:: logger exiting. breaking out on closed log message-queue.\nstarting to flush all the blocked logs.\n")
			close(chanbuffLog)
			for logMsg := range chanbuffLog {
				dumpServerLog(logMsg.logmsg)
			}
			break
		}
	}

	/* for runFlag {
		select {
			case <-doneChan:  // chanbuffLog needs to be closed. pull all the logs from the channel and dump them to file-system.
				runFlag = false
				dumpServerLog("[WARNING]:: logger exiting. breaking out on closed log message-queue.\nstarting to flush all the blocked logs.\n")
				close(chanbuffLog)
				for logMsg := range chanbuffLog {
					dumpServerLog(logMsg.logmsg)
				}
				break
			default:
				break
		}
		select {
			case logMsg, isOK := <-chanbuffLog: // pushes dummy logmessage onto the channel
				if !isOK {
					runFlag = false
					break
				}

				dumpServerLog(logMsg.logmsg)
				break
			default:
				break
		}
	} */
}

/* *****************************************************************************
Description :
- Initializes logger package data.
- Creates a directory $PWD/logs if doesn't exist and creates first logfile
underneath.

Arguments   :
1> isLoggerInit bool: true if logger data to be initialized. false in case logs are sent to stdout and not to any log file.

4> logLevel: either of DEBUG, INFO, WARNING, ERROR.

Return value:
1> bool: True if successful, false otherwise.
*/
func Init(logDir string, logLevel LogLevel) bool {
	var err error

	// check if logger is already initialized
	if isInit {
		return true
	}

	// set log level
	current_LOG_LEVEL = logLevel

	// check if log dir exists
	logDirInfo, err := os.Stat(logDir)
	if err != nil {
		fmt.Printf("Error: Stat(%s): ", logDirInfo, err)
		return false
	}

	logfileNameList = make([]string, log_MAX_FILES)

	chanbuffLog = make(chan logmessage, chanbuffLogCapacity)

	// TODO: check below code
	logFile := filepath.Join(logDir, log_FILE_NAME_PREFIX) + ".1"
	tmplogFile := filepath.Join(logDir, log_FILE_NAME_PREFIX)
	dummyLogfile = logFile + ".dummy"

	pServerLogFile, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("error-4: %s\n", err.Error()) //Error: while creating logfile: %s, error: %s\n", logFile, err.Error())
		return false
	}

	for i := int8(0); i < log_MAX_FILES; i++ {
		logfileNameList[i] = fmt.Sprintf("%s.%d", tmplogFile, i+1)
	}

	errDup2 := syscall.Dup2(int(pServerLogFile.Fd()), syscall.Stdout)
	if errDup2 != nil {
		fmt.Printf("Error: Failed to reuse STDOUT.\n")
	}
	isInit = true
	return true
}
