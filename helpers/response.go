package helpers

import (
	"encoding/json"
	"net/http"
)

const (
	contentType      = "Content-Type"
	contentTypeValue = "application/json; charset=utf-8"

	SUCCESS_RESPONSE    = "Permintaan anda berhasil diproses"
	FAILED_RESPONSE     = "Permintaan anda gagal diproses"
	STATUS_INTERNAL_ERR = "STATUS_INTERNAL_ERROR"
	STATUS_BAD_REQUEST  = "STATUS_BAD_REQUEST"
	STATUS_UNAUTHORIZED = "STATUS_UNAUTHORIZED"
	STATUS_FORBIDDEN    = "STATUS_FORBIDDEN"
	STATUS_NOT_FOUND    = "STATUS_NOT_FOUND"
)

type Response struct {
	Err    interface{} `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Msg    interface{} `json:"message"`
	Status bool        `json:"success"`
}

func SetResponseJSON(data interface{}, message string, status bool) *Response {
	return &Response{
		Data:   data,
		Msg:    message,
		Status: status,
	}
}

func (r *Response) HandleResponse(w http.ResponseWriter, statusCode int) {
	w.Header().Set(contentType, contentTypeValue)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(r)
}

func ErrResponseFieldFormat(errorMessages map[string][]string) map[string][]string {
	errorMap := map[string][]string{}
	for field, messages := range errorMessages {
		errorMap[field] = messages
	}

	return errorMap
}
