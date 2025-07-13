package utils;

import (
		"net/http"
		"encoding/json"
        "fmt"
)

func WriteErrorResponse(w http.ResponseWriter, message string, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
 
    fmt.Printf("Error : %s \n", message)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": message,
        "success": false,
        "error": true,
    })
}

func WriteSuccessResponse(w http.ResponseWriter, message string, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": message,
        "data": data,
        "success": true,
        "error": false,
    })
}