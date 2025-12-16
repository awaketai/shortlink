package idgen

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

const defaultCodeLength = 7

type Generator struct {
}

func NewGenerator() *Generator {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Generator{}
}

func (g *Generator) GenerateShortCode(ctx context.Context, longURL string) (string, error) {
	if longURL == "" {
		return "", fmt.Errorf("idgen: longURL cannot be emptty for code generation")
	}
	hasher := sha256.New()
	hasher.Write([]byte(longURL))
	hasher.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	buf := make([]byte, 0, 20)
	hasher.Write(fmt.Appendf(buf, "%d", rand.Int63()))
	hashBytes := hasher.Sum(nil)
	encoded := base64.URLEncoding.EncodeToString(hashBytes)
	if len(encoded) < defaultCodeLength {
		return encoded, nil
	}

	return encoded[:defaultCodeLength], nil
}
