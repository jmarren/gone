package middleware

import (
	"fmt"
	"net/http"
)

func LogHiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Hi from middleware!\n")
		next.ServeHTTP(w, r)
	})
}

func LogBaconMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Bacon...\n")
		next.ServeHTTP(w, r)
	})
}
