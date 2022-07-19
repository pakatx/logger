package logger

import (
	"fmt"
	"os"
)

/* ****************************************************************************
Description :
- Extracts sourceFilePath - defaultPath from sourceFilePath.

Arguments   :
1> sourceFilePath string: Absolute path of source file where logger.Log() has been called from.
2> defaultPath string: Default path component.

Return value:
1> bool: true is successful, false otherwise.
2> string: Absolute-path less default path.

Additional note: na
**************************************************************************** */
func getFilePath(sourceFilePath string, defaultPath string) (bool, string) {
	filePath := ""
	if len(defaultPath) > len(sourceFilePath) {
		return false, filePath
	}

	fmt.Printf("dbgrm::  sourceFilePath: %s,  defaultPath: %s\n", sourceFilePath, defaultPath)

	length := len(sourceFilePath) - len(defaultPath)
	var i int
	for i = 0; i < length; i++ {
		if sourceFilePath[i] == defaultPath[0] {
			if sourceFilePath[i:i+len(defaultPath)] == defaultPath {
				break
			}
		}
	}

	filePath = sourceFilePath[i+len(defaultPath) : len(sourceFilePath)]
	fmt.Printf("dbgrm::  filePath: %s\n", filePath)
	return true, filePath
}

/* ****************************************************************************
Description :
- Dumps logMsg into target logfile pointed to by plogfile file handler.
- Dumps logMsg into the database table.

Arguments   :
1> logMsg string: log message to be dumped in the logfile.

Return Value: na

Additional note:
TODO: Dump log message into nosql db.
**************************************************************************** */
func dumpServerLog(logMsg string) {
	if pServerLogFile == nil {
		fmt.Printf("error-5\n") // nil file handler
		os.Exit(1)
	}

	if logMsg == "" {
		return
	}

	pServerLogFile.WriteString(logMsg)
	//fmt.Printf(logMsg) // TODO-REM: remove this fmp.Printf() call later

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

		if currentLogfileCnt < 10 {
			currentLogfileCnt = currentLogfileCnt + 1
		}

		go handleLogRotate()
	}
}

/* ****************************************************************************
Description :
- Rotates logs to subsequent log file (n % 10). Each log file is 20MB (20971520 Bytes) size.
- Rolls over and starts from 1st log file if 10th log file is rotated.

Arguments   : na

Return Value: na

Additional note: na
**************************************************************************** */
func handleLogRotate() {
	for i := currentLogfileCnt; i > 2; i-- {
		err := os.Rename(logfileNameList[i-2], logfileNameList[i-1])
		if err != nil {
			// mv %s to %s. error: %s\n", logfileNameList[i-2], logfileNameList[i-1], err.Error())
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
