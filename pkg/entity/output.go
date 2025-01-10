package entity

type Output struct {
	Content        string
	ShowTimeHeader bool
	WithHeader     bool
	CursorAtEnd    bool
}

func (l *Output) DumpLogToChannel(outputChannel chan Output) {
	outputChannel <- *l
}
