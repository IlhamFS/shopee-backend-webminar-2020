package util

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

var (
	jsoni = jsoniter.ConfigCompatibleWithStandardLibrary
)

// ResponseData response data wrapper
type ResponseData struct {
	Data interface{} `json:"data"`
}

// ErrorResponse general error response
type ErrorResponse struct {
	Errors interface{} `json:"errors"`
}

// ErrorDetail specific detail error
type ErrorDetail struct {
	Detail string `json:"detail"`
}

// WriteOKResponse Writes the response as a standard JSON response with StatusOK
func WriteOKResponse(w http.ResponseWriter, m interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	b, err := jsoni.Marshal(&ResponseData{Data: m})
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
	w.Write(b)
}

// WriteErrorResponse Writes the error response as a Standard API JSON response with a response code
func WriteErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)

	b, err := jsoni.Marshal(&ErrorResponse{Errors: &ErrorDetail{Detail: errorMsg}})
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}

	w.Write(b)
}

// WriteCustomResponse Writes the response as a standard JSON response and with a response code
func WriteCustomResponse(w http.ResponseWriter, code int, m interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	b, err := jsoni.Marshal(&ResponseData{Data: m})
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
	w.Write(b)
}
