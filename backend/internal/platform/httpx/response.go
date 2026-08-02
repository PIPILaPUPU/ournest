package httpx

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func ParseID(raw string) (int, error) {
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, err
	}
	return id, nil
}
