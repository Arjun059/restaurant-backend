package utils;

import (
		"net/http"
		"encoding/json"
)

func WriteErrorResponse(w http.ResponseWriter, message string, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": message,
        "success": false,
        "error": true,
    })
}

func WriteSuccessResponse(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "data": data,
        "success": true,
        "error": false,
    })
}