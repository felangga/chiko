package controller

func (c Controller) PrintLog(log string) {
	// Get last log message
	lastLog := c.ui.OutputPanel.GetText(false)
	c.ui.OutputPanel.SetText(lastLog + log)
}
