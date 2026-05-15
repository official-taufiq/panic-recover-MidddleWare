package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", recoverMiddleWare(mux)))
}

func recoverMiddleWare(app http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
				http.Error(w, "Something went wrong:(", http.StatusInternalServerError)
			}
		}()
		app.ServeHTTP(w, r)
	})
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}
