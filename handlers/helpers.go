package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func writeResponse(w http.ResponseWriter, status int, data interface{}, err error) {
	enc := json.NewEncoder(w)
	resp := Response{
		Data: data,
	}
	if err != nil {
		resp.Error = fmt.Sprint(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := enc.Encode(resp); err != nil {
		fmt.Fprintf(w, "error encoding response %v:%s", resp, err)
	}
}
