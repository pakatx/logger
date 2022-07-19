package logger

// log levels
const (
	OFF     LogLevel = iota // 0 - switches off logging
	FATAL                   // 1 - panics
	ERROR                   // 2 - error events
	WARNING                 // 3 - potenially harmful events
	INFO                    // 4 - informational messages
	DEBUG                   // 5 - debug messages - default
)

// buffered channel for logs
const chanbuffLogCapacity int = 10

// log file prefix
const log_FILE_NAME_PREFIX string = "server.log"

// log rotation settings
const log_MAX_FILES uint8 = 10
const log_FILE_SIZE int64 = 20 * 1024 * 1024 // 20 MB
