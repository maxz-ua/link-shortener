package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		var errMsg string
		switch err.ActualTag() {
		case "required":
			errMsg = fmt.Sprintf("field '%s' is required", err.Field())
		case "url":
			errMsg = fmt.Sprintf("field '%s' must be a valid URL", err.Field())
		default:
			errMsg = fmt.Sprintf("field '%s' is not valid", err.Field())
		}
		errMsgs = append(errMsgs, errMsg)
	}

	if len(errMsgs) == 0 {
		return Response{
			Status: StatusOK,
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
