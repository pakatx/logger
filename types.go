package logger

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
		return "invalidLogLevel"
	}
}

type LogMessage struct {
	// common
	TimeStamp string `json:"time_stamp,omitempty"`
	Level     string `json:"level,omitempty"`
	Component string `json:"component,omitempty"`
	// context based
	SourceFile   string `json:"source_file,omitempty"`
	LineNumber   int    `json:"line_number,omitempty"`
	FunctionName string `json:"function_name,omitempty"`
	// webserver
	Method     string `json:"method,omitempty"`
	Path       string `json:"path,omitempty"`
	ClientIP   string `json:"client_ip,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
	Latency    string `json:"latency,omitempty"`
	// common
	Message string `json:"message,omitempty"`
}
