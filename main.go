package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// handleProxyWrapper wraps the anonymous http.HandlerFunc to scope in the targetURL
func handleProxyWrapper(targetURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		println("Received request for:", r.URL.String())

		// This parses the string into a URL struct
		targetURL, err := url.Parse(targetURL)
		if err != nil {
			http.Error(w, "Invalid target URL", http.StatusInternalServerError)
			return
		}

		// Update the target URL's path to match the incoming request's path
		targetURL.Path = r.URL.Path
		targetURL.RawQuery = r.URL.RawQuery

		// Create a new request to forward to the target server
		proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}

		// Copy the headers from the original request
		for k, v := range r.Header {
			proxyReq.Header[k] = v
		}

		// Send the request to the target server
		resp, err := http.DefaultTransport.RoundTrip(proxyReq)
		if err != nil {
			http.Error(w, "Error forwarding request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Copy the response headers and status code
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(resp.StatusCode)

		// Copy the response body
		io.Copy(w, resp.Body)
	}
}

func main() {
	println("Starting proxy server")

	// set up flags
	tgtAddrPtr := flag.String("addr", "localhost", "address to redirect to, shoudl include protocol")
	tgtPortPtr := flag.String("port", "8781", "port to redirect to")
	portPtr := flag.String("lport", "8782", "port to run proxy server on")

	// parse flags
	flag.Parse()

	redirectAddr := fmt.Sprintf("%s:%s", *tgtAddrPtr, *tgtPortPtr)
	serverAddr := fmt.Sprintf(":%s", *portPtr)

	http.HandleFunc("/", handleProxyWrapper(redirectAddr))
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		println("Error starting server:", err)
		return
	}
}
