package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/mitjabez/bite-tracker/views"
)

func main() {
	helloComponent := views.Base("bzy")

	http.Handle("/", templ.Handler(helloComponent))
	http.ListenAndServe(":8000", nil)
}
