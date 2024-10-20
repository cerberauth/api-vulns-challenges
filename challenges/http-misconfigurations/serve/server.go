package serve

import (
	"log"
	"net/http"
	"time"
)

func RunServer(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/http-method-overide", func(w http.ResponseWriter, r *http.Request) {
		method := r.Header.Get("X-HTTP-Method-Override")
		if method == "" {
			method = r.Method
		}

		w.Header().Set("Content-Type", "application/json")
		if method == "POST" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "POST method"}`))
		} else if method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "GET method"}`))
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/headers/cors-wildcard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/headers/csp-frame-ancestors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'http://example.com'")
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/cookies/unsecure", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "unsecure",
			Value:    "unsecure",
			SameSite: http.SameSiteStrictMode,
			Secure:   false,
			HttpOnly: true,
			Expires:  time.Now().Add(24 * time.Hour),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/cookies/not-httponly", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "unsecure",
			Value:    "unsecure",
			SameSite: http.SameSiteStrictMode,
			HttpOnly: false,
			Secure:   true,
			Expires:  time.Now().Add(24 * time.Hour),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/cookies/samesite-none", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "unsecure",
			Value:    "unsecure",
			SameSite: http.SameSiteNoneMode,
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(24 * time.Hour),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/cookies/no-expiration", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "unsecure",
			Value:    "unsecure",
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
			Secure:   true,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
