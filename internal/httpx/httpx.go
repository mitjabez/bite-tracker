package httpx

import (
	"log"
	"net/http"
)

func InternalError(w http.ResponseWriter, msg string, err error) {
	log.Printf("ERROR: %s: %v", msg, err)
	http.Error(w, "Internal Server Error", 500)
}
