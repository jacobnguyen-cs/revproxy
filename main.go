package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

const backendURL = "https://httpbin.org"

var client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

var hopHeaders = []string{"Keep-Alive", "Transfer-Encoding", "Upgrade", "Proxy-Authorization", "Proxy-Authenticate", "Proxy-Connection", "Te", "Trailer"}

func removeHopByHopHeaders(h http.Header) {
	for _, line := range h.Values("Connection") {
		names := strings.Split(line, ",")
		for _, name := range names {
			h.Del(strings.TrimSpace(name))
		}
	}
	h.Del("Connection")

	for _, value := range hopHeaders {
		h.Del(value)
	}
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
	removeHopByHopHeaders(req.Header)

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
	removeHopByHopHeaders(resp.Header)

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
