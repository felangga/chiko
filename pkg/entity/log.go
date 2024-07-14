package entity

type LogType int8

const (
	LOG_INFO    LogType = 0
	LOG_ERROR   LogType = 1
	LOG_WARNING LogType = 2
)

type LogParam struct {
	Content string
	Type    LogType
}