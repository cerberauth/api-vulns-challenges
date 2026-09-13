package serve

import (
	"log"
	"net/http"
)

func RunServer(port string, vulnerable bool) {
	// vulnerable: the request Origin is reflected back with credentials
	// allowed, and the "null" origin is accepted, letting any site (or a
	// sandboxed iframe / file:// page) read authenticated responses
	http.HandleFunc("/api/account", func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Content-Type", "application/json")

		if vulnerable {
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		} else {
			if origin == "https://app.example.com" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Vary", "Origin")
			}
		}

		if r.Method == http.MethodOptions {
			if vulnerable {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "*")
				w.Header().Set("Access-Control-Max-Age", "86400")
			} else {
				w.Header().Set("Access-Control-Allow-Methods", "GET")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Header().Set("Access-Control-Max-Age", "600")
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"email": "user@example.com", "plan": "premium"}`))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
