package text

import (
	"crypto/sha256"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
)

func CreateContentHash(content string) string {
	h := sha256.New()
	contentBytes := []byte(content)
	h.Write(contentBytes)
	contentHash := hex.EncodeToString(h.Sum(nil))
	log.WithField("contentHash", contentHash).Debug("hash")
	return contentHash
}
