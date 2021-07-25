package main

import (
	"net/http"
)

func MethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			method := r.PostFormValue("_method")
			if method == "" {
				method = r.Header.Get("X-HTTP-Method-Override")
			}

			if method == http.MethodPut ||
				method == http.MethodPatch ||
				method == http.MethodDelete {
				r.Method = method
			}
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	println("Hello World!")
}
