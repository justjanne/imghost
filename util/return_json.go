package util

import (
	"encoding/json"
	"net/http"
)

func ReturnJson(writer http.ResponseWriter, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(&data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	} else {
		writer.WriteHeader(200)
		_, err := writer.Write(bytes)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
