package serve

import (
	"log"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// trustedHost is the only virtual host the routing tier is meant to forward
// requests for.
const trustedHost = "app.example.com"

func RunServer(port string, vulnerable bool) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r.Host is populated from the HTTP/2 ":authority" pseudo-header.
		// A legacy "Host" header can also be sent by some HTTP/1.1-era
		// clients or intermediaries during the HTTP/2 downgrade; per RFC
		// 7540 §8.1.2.3 they must match, otherwise the request must be
		// rejected to avoid routing confusion.
		legacyHost := r.Header.Get("Host")

		w.Header().Set("Content-Type", "application/json")

		if !vulnerable {
			if legacyHost != "" && legacyHost != r.Host {
				http.Error(w, `{"error": "authority/host mismatch"}`, http.StatusBadRequest)
				return
			}
		}
		// vulnerable: the mismatch is ignored, the request is routed purely
		// on r.Host (derived from ":authority"), so a crafted legacy Host
		// header that disagrees with ":authority" can smuggle a request past
		// a downstream check that inspects the wrong field
		if r.Host != trustedHost {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "public backend"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "trusted app backend"}`))
	})

	h2s := &http2.Server{}
	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, h2c.NewHandler(handler, h2s)))
}
