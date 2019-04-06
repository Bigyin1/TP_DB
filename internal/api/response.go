package api

import (
	"encoding/json"
	"net/http"
)

func response(rw http.ResponseWriter, code int, payload interface{}) {
	if payload != nil {
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(code)
		json.NewEncoder(rw).Encode(payload)
	} else {
		rw.WriteHeader(code)
	}
}
