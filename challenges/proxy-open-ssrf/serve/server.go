package serve

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

// isPrivate reports whether host resolves to a loopback, link-local or
// private address that must never be reachable from the public fetch proxy.
func isPrivate(host string) bool {
	ips, err := net.LookupIP(host)
	if err != nil {
		return true
	}
	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() {
			return true
		}
	}
	return false
}

// internalSecret simulates a sensitive internal service only meant to be
// reachable from inside the network, never through the public proxy.
func internalSecret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"aws_secret_access_key": "AKIA...leaked"}`))
}

func RunServer(port string, vulnerable bool) {
	http.HandleFunc("/internal/secret", internalSecret)

	// /fetch acts as an open proxy: it fetches whatever URL the caller
	// provides and returns the response body.
	http.HandleFunc("/fetch", func(w http.ResponseWriter, r *http.Request) {
		target := r.URL.Query().Get("url")
		w.Header().Set("Content-Type", "application/json")
		if target == "" {
			http.Error(w, `{"error": "missing url"}`, http.StatusBadRequest)
			return
		}

		parsed, err := url.Parse(target)
		if err != nil {
			http.Error(w, `{"error": "invalid url"}`, http.StatusBadRequest)
			return
		}

		// vulnerable: no allowlist is enforced, so the proxy will happily
		// reach the loopback interface, internal/private IP ranges, or any
		// arbitrary external host on the caller's behalf
		if !vulnerable && isPrivate(parsed.Hostname()) {
			http.Error(w, `{"error": "target host is not allowed"}`, http.StatusForbidden)
			return
		}

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(target)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": %q}`, err.Error()), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, `{"error": "failed to read upstream response"}`, http.StatusBadGateway)
			return
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
