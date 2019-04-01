package api

import (
	"encoding/json"
	"net/http"
)

func response(rw http.ResponseWriter, code int, payload interface{}) {
	if payload != nil {
		rw.WriteHeader(code)
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(payload)
	} else {
		rw.WriteHeader(code)
	}
}
