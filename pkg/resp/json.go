package resp

import (
	"net/http"
	"encoding/json"
)

func Json(w http.ResponseWriter, data any, StatusCode int) {
 w.Header().Set("Content-Type", "application/json")
 w.WriteHeader(StatusCode)
 json.NewEncoder(w).Encode(data)
}