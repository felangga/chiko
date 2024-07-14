package entity

type LogType int8

const (
	LOG_INFO    LogType = 0
	LOG_ERROR   LogType = 1
	LOG_WARNING LogType = 2
)

type Log struct {
	Content string
	Type    LogType
}

func (l *Log) DumpLogToChannel(logChannel chan Log) {
	// Send the log request to channel, then it will be displayed on the log window by the logDumper function
	logChannel <- *l
}
