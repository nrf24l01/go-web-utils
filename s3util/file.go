package s3util

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
)

func FileHashSHA256(fileHeader *multipart.FileHeader) (string, error) {
    file, err := fileHeader.Open()
    if err != nil {
        return "", err
    }
    defer file.Close()

    hasher := sha256.New()
    if _, err := io.Copy(hasher, file); err != nil {
        return "", err
    }

    return hex.EncodeToString(hasher.Sum(nil)), nil
}
