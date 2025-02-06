package entity

type Output struct {
	Content     string
	WithHeader  bool
	CursorAtEnd bool
}

func (l *Output) DumpLogToChannel(outputChannel chan Output) {
	outputChannel <- *l
}
