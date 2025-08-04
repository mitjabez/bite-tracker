package httpx

import (
	"log"
	"net/http"
)

func InternalError(w http.ResponseWriter, msg string, err error) {
	log.Printf("ERROR: %s: %v", msg, err)
	http.Error(w, "Internal Server Error. Try again later.", 500)
}

func BadRequest(w http.ResponseWriter, msg string, err error) {
	log.Printf("WARN: %s: %v", msg, err)
	http.Error(w, "Bad Request", 400)
}
