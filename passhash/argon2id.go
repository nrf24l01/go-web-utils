package passhash

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

	// Expect: ["argon2id", "v=19", "m=...,t=...,p=...", "<salt>", "<hash>"]
	if len(parts) != 5 {
		return false, errors.New("invalid encoded hash format")
	}
	if parts[0] != "argon2id" {
		return false, errors.New("unsupported algorithm")
	}

	// params are in parts[2]
	var memory uint32
	var time uint32
	var parallelism uint8
	_, err := fmt.Sscanf(parts[2], "m=%d,t=%d,p=%d", &memory, &time, &parallelism)
	if err != nil {
		return false, err
	}

	// Try RawStd first, fall back to Std (handles presence/absence of padding)
	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		salt, err = base64.StdEncoding.DecodeString(parts[3])
		if err != nil {
			return false, err
		}
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		hash, err = base64.StdEncoding.DecodeString(parts[4])
		if err != nil {
			return false, err
		}
	}

	computedHash := argon2.IDKey([]byte(password), salt, time, memory, parallelism, uint32(len(hash)))

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
