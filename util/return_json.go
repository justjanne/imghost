package util

import (
	"encoding/json"
	"net/http"
)

func ReturnJson(writer http.ResponseWriter, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(data); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
