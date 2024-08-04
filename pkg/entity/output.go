package entity

type Output struct {
	Content        string
	ShowTimeHeader bool
}

func (l *Output) DumpLogToChannel(outputChannel chan Output) {
	outputChannel <- *l
}
