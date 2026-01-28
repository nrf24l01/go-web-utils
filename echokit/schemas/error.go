package schemas

import (
	"time"

	"github.com/labstack/echo/v4"
)

type ErrorCode string

const (
	BAD_REQUEST              ErrorCode = "BAD_REQUEST"
	VALIDATION_FAILED        ErrorCode = "VALIDATION_FAILED"
	UNAUTHORIZED             ErrorCode = "UNAUTHORIZED"
	FORBIDDEN                ErrorCode = "FORBIDDEN"
	NOT_FOUND                ErrorCode = "NOT_FOUND"
	USER_NOT_FOUND           ErrorCode = "USER_NOT_FOUND"
	EMAIL_ALREADY_EXISTS     ErrorCode = "EMAIL_ALREADY_EXISTS"
	INTERNAL_SERVER_ERROR    ErrorCode = "INTERNAL_SERVER_ERROR"
)

type ApiError struct {
	Code      ErrorCode              `json:"code"`
	Message   string                 `json:"message"`
	TraceID   string                 `json:"traceId"`
	Timestamp time.Time              `json:"timestamp"`
	Path      string                 `json:"path"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

type FieldError struct {
	Field         string      `json:"field"`
	Issue         string      `json:"issue"`
	RejectedValue interface{} `json:"rejectedValue,omitempty"`
}

type ValidationError struct {
	ApiError
	FieldErrors []FieldError `json:"fieldErrors"`
}

type DslError struct {
	Code     string  `json:"code"`
	Message  string  `json:"message"`
	Position *int    `json:"position,omitempty"`
	Near     *string `json:"near,omitempty"`
}

func GenError(c echo.Context, code ErrorCode, message string, details map[string]interface{}) ApiError {
	traceID := ""
	if v := c.Get("traceId"); v != nil {
		if s, ok := v.(string); ok {
			traceID = s
		}
	}

	ts := time.Now()
	if v := c.Get("timestamp"); v != nil {
		switch t := v.(type) {
		case time.Time:
			ts = t
		case string:
			if parsed, err := time.Parse(time.RFC3339, t); err == nil {
				ts = parsed
			}
		}
	}

	res := ApiError{
		Code:      code,
		Message:   message,
		TraceID:   traceID,
		Timestamp: ts,
		Path:      c.Request().URL.Path,
		Details:   details,
	}
	return res
}


// CustomErrorCode converts a plain string into the typed ErrorCode.
// Use when you want an ErrorCode value from a dynamic string.
func CustomErrorCode(s string) ErrorCode {
	return ErrorCode(s)
}
