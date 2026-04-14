package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func Send(w http.ResponseWriter, status int, data any, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := APIResponse{
		Data:  data,
		Error: errMsg,
	}

	json.NewEncoder(w).Encode(res)
}
