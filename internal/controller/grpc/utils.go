package grpc

import (
	"regexp"
)

var re = regexp.MustCompile(`\(([^()]*)\)`)

func (g *GRPC) parseRequestResponse(text string) [][]string {
	return re.FindAllStringSubmatch(text, -1)
}
