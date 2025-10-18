package validators

import (
	"errors"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/go-playground/validator/v10"
)

// ParseSize parses human-readable size like "5MB", "10KB", etc.
func ParseSize(param string) (int64, error) {
	if param == "" {
		return 0, errors.New("size parameter empty")
	}

	param = strings.TrimSpace(strings.ToUpper(param))
	multiplier := int64(1)
	switch {
	case strings.HasSuffix(param, "KB"):
		multiplier = 1 << 10
		param = strings.TrimSuffix(param, "KB")
	case strings.HasSuffix(param, "MB"):
		multiplier = 1 << 20
		param = strings.TrimSuffix(param, "MB")
	case strings.HasSuffix(param, "GB"):
		multiplier = 1 << 30
		param = strings.TrimSuffix(param, "GB")
	case strings.HasSuffix(param, "B"):
		param = strings.TrimSuffix(param, "B")
	}

	value, err := strconv.ParseInt(strings.TrimSpace(param), 10, 64)
	if err != nil {
		return 0, err
	}
	return value * multiplier, nil
}

// filetype validation tag handler
// Usage: `validate:"filetype=image/png|image/jpeg"`
func fileTypeValidator(fl validator.FieldLevel) bool {
	field := fl.Field().Interface()
	allowedParam := fl.Param()
	if allowedParam == "" {
		return true // no restriction
	}

	allowedTypes := strings.Split(allowedParam, "|")
	allowedSet := make(map[string]struct{}, len(allowedTypes))
	for _, t := range allowedTypes {
		allowedSet[strings.TrimSpace(t)] = struct{}{}
	}

	switch v := field.(type) {
	case *multipart.FileHeader:
		return validateSingleFileType(v, allowedSet)
	case []*multipart.FileHeader:
		for _, f := range v {
			if !validateSingleFileType(f, allowedSet) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func validateSingleFileType(fh *multipart.FileHeader, allowed map[string]struct{}) bool {
	if fh == nil {
		return false
	}

	src, err := fh.Open()
	if err != nil {
		return false
	}
	defer src.Close()

	mtype, err := mimetype.DetectReader(src)
	if err != nil {
		return false
	}

	if _, ok := allowed[mtype.String()]; !ok {
		return false
	}
	return true
}

// filesize validation tag handler
// Usage: `validate:"filesize=5MB"`
func fileSizeValidator(fl validator.FieldLevel) bool {
	field := fl.Field().Interface()
	param := fl.Param()
	if param == "" {
		return true
	}

	maxBytes, err := ParseSize(param)
	if err != nil {
		return false
	}

	switch v := field.(type) {
	case *multipart.FileHeader:
		return v != nil && v.Size <= maxBytes
	case []*multipart.FileHeader:
		for _, f := range v {
			if f == nil || f.Size > maxBytes {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// RegisterFileValidations adds `filetype` and `filesize` validators to a validator instance.
func RegisterFileValidations(v *validator.Validate) error {
	if err := v.RegisterValidation("filetype", fileTypeValidator); err != nil {
		return err
	}
	if err := v.RegisterValidation("filesize", fileSizeValidator); err != nil {
		return err
	}
	return nil
}
