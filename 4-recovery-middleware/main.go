package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", recoveryHandler(mux)))
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func recoveryHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w}
		defer func() {
			if err := recover(); err != nil {
				stackTrace := debug.Stack()
				logrus.WithFields(logrus.Fields{
					"error": err,
					"stack": string(stackTrace),
				}).Error("panic occurred")

				if !rw.wroteHeader {
					rw.WriteHeader(http.StatusInternalServerError)
				}

				if os.Getenv("ENV") == "development" {
					fmt.Fprintf(rw, "panic: %v\n\nStack Trace:\n%s", err, stackTrace)
				} else {
					fmt.Fprint(rw, "Something went wrong")
				}
			}
		}()
		next.ServeHTTP(rw, r)
	})
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
