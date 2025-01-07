package entity

type Output struct {
	Content        string
	ShowTimeHeader bool
	WithHeader     bool
}

func (l *Output) DumpLogToChannel(outputChannel chan Output) {
	outputChannel <- *l
}
