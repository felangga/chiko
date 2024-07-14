package controller

import (
	"regexp"
)

func (c *Controller) parseRequestResponse(text string) [][]string {
	var re = regexp.MustCompile(`\(([^()]*)\)`)
	return re.FindAllStringSubmatch(text, -1)
}
