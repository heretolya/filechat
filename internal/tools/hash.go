package tools

import (
	"crypto/sha256"
	"fmt"
)

func Hash(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	bs := h.Sum(nil)
	hash := fmt.Sprintf("%x", bs)
	return hash
}
