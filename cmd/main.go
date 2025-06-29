package main

import (
	"github.com/a-h/templ"
	"github.com/mitjabez/bite-tracker/views"
	"net/http"
)

func main() {
	helloComponent := views.Base(views.Hello("bzy"), "Hello")
	helpComponent := views.Base(views.Help(), "Help")
	assetsHandler := http.FileServer(http.Dir("views/assets"))

	http.Handle("/", templ.Handler(helloComponent))
	http.Handle("/help", templ.Handler(helpComponent))
	http.Handle("/assets/", http.StripPrefix("/assets", assetsHandler))
	http.ListenAndServe(":8000", nil)
}
