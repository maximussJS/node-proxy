package responses

import (
	"log"
	"net/http"
)

func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte("Internal server error:"))
	if err != nil {
		log.Fatalf("Error writing internal server error response %s", err)
	}
}

func BadRequest(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, e := w.Write([]byte(err.Error()))
	if e != nil {
		log.Fatalf("Error writing bad request response %s", e)
	}
}

func RequestTimeout(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusRequestTimeout)
	_, err := w.Write([]byte(`{"error": "Request timeout"}`))
	if err != nil {
		log.Fatalf("Error writing request timeout response %s", err)
	}
}
