package validators

import (
	"errors"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/go-playground/validator/v10"
)

// ParseSize parses human-readable size strings like "5MB", "10KB", "512B".
// Returns size in bytes or error.
func ParseSize(param string) (int64, error) {
	if param == "" {
		return 0, errors.New("empty size parameter")
	}
	s := strings.TrimSpace(strings.ToUpper(param))

	mult := int64(1)
	switch {
	case strings.HasSuffix(s, "KB"):
		mult = 1 << 10
		s = strings.TrimSuffix(s, "KB")
	case strings.HasSuffix(s, "MB"):
		mult = 1 << 20
		s = strings.TrimSuffix(s, "MB")
	case strings.HasSuffix(s, "GB"):
		mult = 1 << 30
		s = strings.TrimSuffix(s, "GB")
	case strings.HasSuffix(s, "B"):
		mult = 1
		s = strings.TrimSuffix(s, "B")
	}

	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("invalid size parameter")
	}
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return val * mult, nil
}

// splitAllowedParam splits a parameter that lists allowed MIME types.
// It accepts ',' or ';' as separators. Trims whitespace and drops empties.
func splitAllowedParam(param string) []string {
	if param == "" {
		return nil
	}
	var parts []string
	// Prefer ';' if present, otherwise comma, otherwise treat whole as single.
	if strings.Contains(param, ";") {
		parts = strings.Split(param, ";")
	} else if strings.Contains(param, ",") {
		parts = strings.Split(param, ",")
	} else {
		parts = []string{param}
	}
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// validateSingleFileType checks MIME of the provided FileHeader using mimetype.DetectReader.
// Returns true if detected mime is in allowed set.
func validateSingleFileType(fh *multipart.FileHeader, allowed map[string]struct{}) bool {
	if fh == nil {
		return false
	}
	src, err := fh.Open()
	if err != nil {
		return false
	}
	defer src.Close()

	// DetectReader can read more than the first 512 bytes, making detection more reliable.
	mtype, err := mimetype.DetectReader(src)
	if err != nil {
		return false
	}
	if _, ok := allowed[mtype.String()]; !ok {
		return false
	}
	return true
}

// fileTypeValidator is the validator tag handler for "filetype".
// Usage in struct tag: `validate:"filetype=image/png,image/jpeg"`
// NOTE: do NOT use '|' inside the parameter (validator internal parsing will split on |).
func fileTypeValidator(fl validator.FieldLevel) bool {
	param := fl.Param()
	if strings.TrimSpace(param) == "" {
		// No restriction declared -> consider valid.
		return true
	}

	mimes := splitAllowedParam(param)
	if len(mimes) == 0 {
		// Malformed param -> treat as invalid configuration -> fail validation
		return false
	}
	allowed := make(map[string]struct{}, len(mimes))
	for _, m := range mimes {
		allowed[m] = struct{}{}
	}

	field := fl.Field().Interface()
	switch v := field.(type) {
	case *multipart.FileHeader:
		return validateSingleFileType(v, allowed)
	case []*multipart.FileHeader:
		for _, fh := range v {
			if !validateSingleFileType(fh, allowed) {
				return false
			}
		}
		return true
	default:
		// Unsupported field type
		return false
	}
}

// fileSizeValidator is the validator tag handler for "filesize".
// Usage: `validate:"filesize=5MB"`
func fileSizeValidator(fl validator.FieldLevel) bool {
	param := strings.TrimSpace(fl.Param())
	if param == "" {
		// no size limit set -> valid
		return true
	}
	maxBytes, err := ParseSize(param)
	if err != nil {
		// malformed size param -> fail validation
		return false
	}

	field := fl.Field().Interface()
	switch v := field.(type) {
	case *multipart.FileHeader:
		if v == nil {
			return false
		}
		return v.Size <= maxBytes
	case []*multipart.FileHeader:
		for _, fh := range v {
			if fh == nil || fh.Size > maxBytes {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// RegisterFileValidations registers the "filetype" and "filesize" validators
// into a provided *validator.Validate instance.
func RegisterFileValidations(v *validator.Validate) error {
	if v == nil {
		return errors.New("validator instance is nil")
	}
	if err := v.RegisterValidation("filetype", fileTypeValidator); err != nil {
		return err
	}
	if err := v.RegisterValidation("filesize", fileSizeValidator); err != nil {
		return err
	}
	return nil
}