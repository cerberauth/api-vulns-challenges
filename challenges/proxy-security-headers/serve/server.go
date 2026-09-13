package serve

import (
	"log"
	"net/http"
)

func RunServer(port string, vulnerable bool) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// vulnerable: the proxy forwards the backend response as-is, without
		// injecting any of the common hardening headers, and leaks the
		// upstream software banner
		if vulnerable {
			w.Header().Set("Server", "nginx/1.18.0 (Ubuntu)")
			w.Header().Set("X-Powered-By", "Express")
		} else {
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("Referrer-Policy", "no-referrer")
			w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "ok"}`))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
