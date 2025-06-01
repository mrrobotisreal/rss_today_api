package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func (app *models.App) GenerateContentHash(title, link string) string {
	hash := sha256.Sum256([]byte(title + link))
	return hex.EncodeToString(hash[:])
}
