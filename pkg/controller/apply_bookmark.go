package controller

import (
	"chiko/pkg/entity"
)

func (c *Controller) ApplyBookmark(session entity.Session) {
	// Get selected connection
	*c.Conn = session

	go func() {
		err := c.CheckGRPC(c.Conn.ServerURL)
		if err != nil {
			c.PrintLog(entity.LogParam{
				Content: err.Error(),
				Type:    entity.LOG_ERROR,
			})
			return
		}
	}()
}
