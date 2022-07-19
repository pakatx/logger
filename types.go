package logger

type LogLevel int

type logmessage struct {
	componentFlag int8
	component     string
	logmsg        string
}

type LogConfig struct {
	SrcBaseDir      string `json:"srcBaseDir"`      // $PWD
	FileSize        int    `json:"fileSize"`        // 20971520 (20MB)
	MaxFilesCnt     int    `json:"maxFilesCnt"`     // 10
	DefaultLogLevel string `json:"defaultLogLevel"` // DEBUG
}
