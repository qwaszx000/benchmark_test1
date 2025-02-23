package main

import (
	"net/http"
)

func handle_test_endpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || r.URL.Path != "/test_plain" {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("Hello world!"))
}

func main() {
	http.HandleFunc("/test_plain", handle_test_endpoint)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
