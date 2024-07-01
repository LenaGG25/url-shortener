package service

import (
	"crypto/sha256"
	"encoding/base64"
)

type URLGenerator struct {
	maxShortKeySize int
}

func NewURLGenerator(maxShortKeySize int) *URLGenerator {
	return &URLGenerator{
		maxShortKeySize: maxShortKeySize,
	}
}

// GenerateShortURL from original url using sha256 hash.
// In the future, original url can be concatenated with user id.
func (g *URLGenerator) GenerateShortURL(originalURL string) string {
	hasher := sha256.New()

	hasher.Write([]byte(originalURL))
	hashed := hasher.Sum(nil)

	shortURL := base64.URLEncoding.EncodeToString(hashed)
	return shortURL[:g.maxShortKeySize]
}
