package serve

import (
	"fmt"
	"log"
	"net/http"
)

const canonicalHost = "app.example.com"

func RunServer(port string, vulnerable bool) {
	// vulnerable: the password-reset link is built from the client-supplied
	// Host header, so an attacker can poison it to point at a host they
	// control and hijack the reset token
	http.HandleFunc("/password-reset", func(w http.ResponseWriter, r *http.Request) {
		host := canonicalHost
		if vulnerable {
			host = r.Host
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"resetLink": "https://%s/reset?token=abc123"}`, host)
	})

	// vulnerable: a crafted Host header changes which backend the request is
	// routed to, allowing virtual-host confusion
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if vulnerable && r.Host == "internal-admin.local" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "internal admin backend"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "public backend"}`))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
