package grpc

import (
	"regexp"
)

func (g *GRPC) parseRequestResponse(text string) [][]string {
	var re = regexp.MustCompile(`\(([^()]*)\)`)
	return re.FindAllStringSubmatch(text, -1)
}
