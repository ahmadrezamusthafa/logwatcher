package errors

import (
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"net/http"
	"runtime"
	"strings"
)

type (
	Error struct {
		ID         int    `json:"-"`
		Code       int    `json:"code"`
		Title      string `json:"title"`
		Detail     string `json:"detail,omitempty"`
		DetailId   string `json:"detail_id"`
		DetailEn   string `json:"detail_en"`
		Messages   string `json:"messages,omitempty"`
		HttpStatus int    `json:"http_status,omitempty"`
		ErrorWithDetails
	}

	ErrorWithDetails struct {
		Err            error    `json:"-"`
		TrueErrMessage string   `json:"-"`
		Traces         []string `json:"traces"`
	}
)

func (e Error) Error() string {
	buf, err := jsoniter.Marshal(e)
	if err != nil {
		return ""
	}
	return string(buf)
}

func (h ErrorWithDetails) Error() string {
	return h.Err.Error()
}

func New(text string) error {
	return errors.New(text)
}

func AddTrace(err interface{}) error {

	if err == nil {
		return nil
	}

	if mainErr, ok := err.(Error); ok {
		mainErr.Traces = append(mainErr.Traces, getLineOfCode(2))
		return mainErr
	}

	if detailErr, ok := err.(ErrorWithDetails); ok {
		detailErr.Traces = append(detailErr.Traces, getLineOfCode(2))
		return detailErr
	}

	parentErr := Error{
		Code:       900001,
		Title:      "Default Error",
		DetailEn:   err.(error).Error(),
		HttpStatus: http.StatusBadRequest,
		ErrorWithDetails: ErrorWithDetails{
			Err:    err.(error),
			Traces: []string{getLineOfCode(2)},
		},
	}

	return parentErr
}

func GetTrace(err interface{}) []string {
	if _, ok := err.(error); !ok {
		return []string{}
	}

	if mainErr, ok := err.(Error); ok {
		return mainErr.Traces
	}

	if detailErr, ok := err.(ErrorWithDetails); ok {
		return detailErr.Traces
	}

	return []string{}
}

func ParseError(lang string, err error) Error {

	mainErr, ok := err.(Error)
	if !ok {
		return Error{}
	}

	if lang == "en" {
		mainErr.Detail = mainErr.DetailEn
	} else {
		mainErr.Detail = mainErr.DetailId
	}

	return mainErr
}

func GetHttpStatus(err interface{}) int {
	if mainErr, ok := err.(Error); ok {
		return mainErr.HttpStatus
	}
	return http.StatusInternalServerError
}

func getLineOfCode(skip int) string {
	_, filePath, line, _ := runtime.Caller(skip)
	path := strings.Split(filePath, "/")
	if len(path) > 4 {
		start := len(path) - 4
		filePath = "/" + strings.Join(path[start:], "/")
	}

	details := fmt.Sprintf("%s[%d]", filePath, line)
	return details
}

func Match(firstErr, secondErr error) bool {
	if firstErr == nil && secondErr == nil {
		return true
	}

	var firstCode, secondCode int
	if firstErr != nil && secondErr != nil {
		err1, ok := firstErr.(*Error)
		if ok {
			firstCode = err1.Code
		}

		err2, ok := secondErr.(*Error)
		if ok {
			secondCode = err2.Code
		}

		if firstCode == secondCode {
			return true
		}
	}

	return false
}

var (
	GeneralError                error = Error{Code: 100002, Title: "General Error", DetailEn: "Something went wrong please try again", DetailId: "Mohon maaf permohonan Anda tidak dapat diproses, mohon coba kembali", HttpStatus: http.StatusBadRequest}
	MarshalError                error = Error{Code: 200001, Title: "Library Error", DetailEn: "Cannot marshall data", HttpStatus: http.StatusBadRequest}
	UnmarshalError              error = Error{Code: 200002, Title: "Library Error", DetailEn: "Cannot unmarshall data", HttpStatus: http.StatusBadRequest}
	ReadDataError               error = Error{Code: 200003, Title: "Library Error", DetailEn: "Cannot read data", HttpStatus: http.StatusNotFound}
	AlreadyExistsError          error = Error{Code: 300001, Title: "Data Error", DetailEn: "Data is already exist", HttpStatus: http.StatusPreconditionFailed}
	SqlNoRowsError              error = Error{Code: 300002, Title: "Data Error", DetailEn: "No rows selected", HttpStatus: http.StatusNotFound}
	InvalidUsernameError        error = Error{Code: 400001, Title: "Request Error", DetailEn: "Invalid username", HttpStatus: http.StatusBadRequest}
	InvalidPasswordError        error = Error{Code: 400002, Title: "Request Error", DetailEn: "Invalid user password", HttpStatus: http.StatusBadRequest}
	InvalidRequestParamError    error = Error{Code: 400003, Title: "Request Error", DetailEn: "Invalid request parameter", HttpStatus: http.StatusPreconditionFailed}
	IncompleteRequestParamError error = Error{Code: 400004, Title: "Request Error", DetailEn: "Incomplete request parameter", HttpStatus: http.StatusPreconditionFailed}
)
