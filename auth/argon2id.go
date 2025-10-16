package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Params struct {
	Memory      uint32
	Time        uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var DefaultParams = &Params{
	Memory:      64 * 1024, // 64 MB
	Time:        3,
	Parallelism: 4,
	SaltLength:  16,
	KeyLength:   32,
}

func generateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func HashPassword(password string, p *Params) (string, error) {
	salt, err := generateSalt(p.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Parallelism, p.KeyLength)

	// Возвращаем всё одной строкой
	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		p.Memory, p.Time, p.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash))

	return encoded, nil
}

func CheckPassword(password string, encodedHash string) (bool, error) {
	encodedHash = strings.TrimSpace(encodedHash)

	// Split and remove empty segments so leading/trailing '$' don't break parsing
	rawParts := strings.Split(encodedHash, "$")
	parts := make([]string, 0, len(rawParts))
	for _, p := range rawParts {
		if p != "" {
			parts = append(parts, p)
		}
	}

	// Ensure algorithm present
	foundAlg := false
	for _, p := range parts {
		if p == "argon2id" {
			foundAlg = true
			break
		}
	}
	if !foundAlg {
		return false, errors.New("unsupported algorithm or invalid encoded hash format")
	}

	// Need at least params + salt + hash
	if len(parts) < 3 {
		return false, errors.New("invalid encoded hash format")
	}

	// Find params part (one that contains m=)
	var paramsPart string
	for _, p := range parts {
		if strings.Contains(p, "m=") && strings.Contains(p, "t=") && strings.Contains(p, "p=") {
			paramsPart = p
			break
		}
	}
	if paramsPart == "" {
		return false, errors.New("parameters not found in encoded hash")
	}

	// Salt and hash are expected to be the last two parts
	if len(parts) < 2 {
		return false, errors.New("invalid encoded hash format")
	}
	saltB64 := parts[len(parts)-2]
	hashB64 := parts[len(parts)-1]

	// Parse params
	var memory, time uint32
	var parallelism uint32
	_, err := fmt.Sscanf(paramsPart, "m=%d,t=%d,p=%d", &memory, &time, &parallelism)
	if err != nil {
		return false, err
	}

	// Try RawStd first, fall back to Std (handles presence/absence of padding)
	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		salt, err = base64.StdEncoding.DecodeString(saltB64)
		if err != nil {
			return false, err
		}
	}

	hash, err := base64.RawStdEncoding.DecodeString(hashB64)
	if err != nil {
		hash, err = base64.StdEncoding.DecodeString(hashB64)
		if err != nil {
			return false, err
		}
	}

	if len(hash) == 0 {
		return false, errors.New("decoded hash is empty")
	}

	computedHash := argon2.IDKey([]byte(password), salt, time, memory, uint8(parallelism), uint32(len(hash)))

	return subtleCompare(hash, computedHash), nil
}

func subtleCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
