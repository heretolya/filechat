package tools

import (
	"regexp"
	"strings"
)

func Tokenize(text string) []string {
	text = strings.ToLower(text)
	re := regexp.MustCompile(`[^a-z0-9 ]+`)
	text = re.ReplaceAllString(text, " ")
	return strings.Fields(text)
}
