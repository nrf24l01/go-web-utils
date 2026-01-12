package misc

import (
	"crypto/rand"
	"encoding/binary"
	mrand "math/rand"
)

const charset string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type RandomGenerator struct {
	mr *mrand.Rand
	charset string
}

type RandomGeneratorConfig struct {
	Charset string
}

func NewRandomGenerator(cfg RandomGeneratorConfig) (RandomGenerator, error) {
	var seed int64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &seed); err != nil {
		return RandomGenerator{}, err
	}
	if cfg.Charset != "" {
		return RandomGenerator{
			mr:      mrand.New(mrand.NewSource(seed)),
			charset: cfg.Charset,
		}, nil
	}
	return RandomGenerator{
		mr:      mrand.New(mrand.NewSource(seed)),
		charset: charset,
	}, nil
}

func (rg *RandomGenerator) RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = rg.charset[rg.mr.Intn(len(rg.charset))]
	}
	return string(b)
}

func (rg *RandomGenerator) RandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (rg *RandomGenerator) RandomInt(min, max int) int {
	return rg.mr.Intn(max-min) + min
}