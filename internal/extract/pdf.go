package extract

import (
	"fmt"
	"strings"

	"filechat/pkg/poppler"
)

func PDF(file string) (string, error) {
	doc, err := poppler.Open(file)
	if err != nil {
		return "", err
	}
	defer doc.Close()
	var buf strings.Builder
	for i := range doc.GetNPages() {
		page := doc.GetPage(i)
		text := page.Text()
		fmt.Fprintln(&buf, text)
		page.Close()
	}
	result := buf.String()
	return result, nil
}
