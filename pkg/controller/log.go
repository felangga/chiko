package controller

import "chiko/pkg/entity"

func (c *Controller) PrintLog(param entity.LogParam) {
	// Send the log request to channel, then it will be displayed on the log window by the logDumper function
	c.LogDump <- param
}
