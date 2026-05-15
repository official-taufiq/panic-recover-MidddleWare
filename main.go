package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", recoverMiddleWare(mux, false)))
}

func recoverMiddleWare(app http.Handler, dev bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Println(rec)
				stack := debug.Stack()
				if !dev {
					http.Error(w, "Something went wrong:(", http.StatusInternalServerError)
					return
				}
				fmt.Fprintf(w, "<h1>Panic: %v</h1><pre>%s</pre>", rec, string(stack))
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
