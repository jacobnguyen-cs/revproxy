package main

import (
	"io"
	"log"
	"net/http"
)

const backendURL = "https://httpbin.org"

var client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequestWithContext(r.Context(), r.Method, backendURL+r.RequestURI, r.Body)
	if err != nil {
		log.Printf("create request error: %v", err)
		http.Error(w, "Bad Gateway", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header.Clone()
	req.ContentLength = r.ContentLength

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("backend error: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("could not copy response body: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", proxyHandler)

	log.Print("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
