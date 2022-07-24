package logger

import (
	"fmt"
	"os"
)

// dumpServerLog dumps logMessage into local files and forwards to log aggregator
// calls handleLogRotate if size of log file exceeds limit.
func dumpServerLog(logMessage LogMessage) {
	if pServerLogFile == nil {
		logMessage.Message = "dumpServerLog Fatal: Corrupt log file handler; Original Message: " + logMessage.Message
		logMessage.Level = "FATAL"
	}

	// Forward to Log Aggregator
	// jsonMessage, _ := json.Marshal(logMessage)
	// jsonMessageStr := string(jsonMessage[:])
	// pServerLogFile.WriteString(jsonMessageStr)
	// TODO: Sameer Oak: Invoke Log Aggregator API here

	pServerLogFile.WriteString(logMessage.String())

	fi, err := pServerLogFile.Stat()
	if err != nil {
		fmt.Printf("error-6: %s\n", err.Error()) // Couldn't obtain stat
		return
	}

	fileSize := fi.Size()
	if fileSize >= log_FILE_SIZE {
		pServerLogFile.Close()
		pServerLogFile = nil
		err = os.Rename(logfileNameList[0], dummyLogfile)
		if err != nil {
			fmt.Printf("error-7: %s\n", err.Error()) // mv %s to %s, error: %s\n", logfileNameList[0], dummyLogfile, err.Error())
			pServerLogFile, err = os.OpenFile(logfileNameList[0], os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
			return
		}

		pServerLogFile, err = os.OpenFile(logfileNameList[0], os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("error-8: %s\n", err.Error()) // recreating logfile: %s,  error: %s\n", logfileNameList[0], err.Error())
			return
		}

		if currentLogfileCnt < log_MAX_FILES {
			currentLogfileCnt = currentLogfileCnt + 1
		}

		go handleLogRotate()
	}
}

// handleLogRotate rotates logs to subsequent log file (n % log_MAX_FILES).
// Each log file is log_FILE_SIZE
// Rolls over and starts from 1st log file if log_MAX_FILES log file is rotated.
func handleLogRotate() {
	for i := currentLogfileCnt; i > 2; i-- {
		err := os.Rename(logfileNameList[i-2], logfileNameList[i-1])
		if err != nil {
			fmt.Printf("error-10: %s\n", err.Error())
			return
		}
	}
	err := os.Rename(dummyLogfile, logfileNameList[1])
	if err != nil {
		// while mv %s to %s. error: %s\n", dummyLogfile, logfileNameList[1], err.Error())
		fmt.Printf("error-11: %s\n", err.Error())
		return
	}
}
