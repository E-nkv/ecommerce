package handlers

import (
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, st int, msg string) {
	http.Error(w, msg, st)
}

func writeJSON(w http.ResponseWriter, st int, data any) {
	bs, err := json.Marshal(data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error marshaling response")
		return
	}
	w.WriteHeader(st)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}
