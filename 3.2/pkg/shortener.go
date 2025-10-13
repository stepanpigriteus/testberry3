package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func Shortener(link string) string {
	h := sha256.Sum256([]byte(link + time.Now().Format("2025-10-10 18:23:45")))
	return hex.EncodeToString(h[:])[:6]
}
