package serve

import (
	"log"
	"net/http"
	"time"
)

func RunServer(port string, vulnerable bool) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/http-method-override", func(w http.ResponseWriter, r *http.Request) {
		validToken := "valid-token"
		if r.Header.Get("Authorization") != "Bearer "+validToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// vulnerable: a GET-only endpoint can also be reached with the real
		// method overridden via a header or query parameter, which lets an
		// attacker bypass method-based access controls (e.g. a proxy/WAF
		// rule that only inspects r.Method)
		overridden := vulnerable && (r.Header.Get("X-HTTP-Method-Override") == http.MethodGet || r.URL.Query().Get("_method") == http.MethodGet)
		if r.Method == http.MethodGet || overridden {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "GET method"}`))
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/headers/cors-wildcard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if vulnerable {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "https://trusted.example.com")
			w.Header().Set("Vary", "Origin")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/headers/csp-frame-ancestors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if vulnerable {
			w.Header().Set("Content-Security-Policy", "frame-ancestors 'http://example.com'")
		} else {
			w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/cookies/unsecure", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "unsecure",
			Value:    "unsecure",
			SameSite: http.SameSiteStrictMode,
			Secure:   !vulnerable,
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
			HttpOnly: !vulnerable,
			Secure:   true,
			Expires:  time.Now().Add(24 * time.Hour),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/cookies/samesite-none", func(w http.ResponseWriter, r *http.Request) {
		sameSite := http.SameSiteNoneMode
		if !vulnerable {
			sameSite = http.SameSiteStrictMode
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "unsecure",
			Value:    "unsecure",
			SameSite: sameSite,
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(24 * time.Hour),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/cookies/no-expiration", func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:     "unsecure",
			Value:    "unsecure",
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
			Secure:   true,
		}
		if !vulnerable {
			cookie.Expires = time.Now().Add(24 * time.Hour)
		}
		http.SetCookie(w, cookie)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
