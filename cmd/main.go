package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/mitjabez/bite-tracker/views"
)

func main() {
	helloComponent := views.Base(views.Hello("bzy"), "Hello")
	helpComponent := views.Base(views.Help(), "Help")

	http.Handle("/", templ.Handler(helloComponent))
	http.Handle("/help", templ.Handler(helpComponent))
	http.ListenAndServe(":8000", nil)
}
