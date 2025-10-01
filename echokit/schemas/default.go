package schemas

import "net/http"

var DefaultInternalErrorResponse = ErrorResponse{
	Message: "Internal Server Error",
	Code:    http.StatusInternalServerError,
}

var DefaultPermissionDeniedResponse = ErrorResponse{
	Message: "Permission Denied",
	Code:    http.StatusForbidden,
}

var DefaultBadRequestResponse = ErrorResponse{
	Message: "Bad Request",
	Code:    http.StatusBadRequest,
}

var DefaultNotFoundResponse = ErrorResponse{
	Message: "Not Found",
	Code:    http.StatusNotFound,
}