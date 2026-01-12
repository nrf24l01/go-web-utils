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

func newFastRand(cfg RandomGeneratorConfig) (RandomGenerator, error) {
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