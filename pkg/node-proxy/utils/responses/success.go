package responses

import (
	"log"
	"net/http"
)

func Success(w http.ResponseWriter, responseBody []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(responseBody)
	if err != nil {
		log.Fatalln("Error writing success response", err)
	}
}
