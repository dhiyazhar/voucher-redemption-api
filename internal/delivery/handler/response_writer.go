package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dhiyazhar/voucher-redemption-api/internal/model"
)

func JSON(w http.ResponseWriter, code int, status string, data any, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := model.WebResponse[any]{
		Code:   code,
		Status: status,
		Data:   data,
		Errors: err,
	}

	json.NewEncoder(w).Encode(resp)
}

func Success(w http.ResponseWriter, code int, status string, data any) {
	JSON(w, code, status, data, "")
}

func Error(w http.ResponseWriter, code int, status string, err string) {
	JSON(w, code, status, nil, err)
}
