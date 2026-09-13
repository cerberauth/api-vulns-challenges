package serve

import (
	"fmt"
	"log"
	"net/http"
)

func RunServer(port string, vulnerable bool) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "ok"}`))
	})

	// vulnerable: a verbose error page leaks the backend's internal path,
	// hostname and a fake stack trace to the client
	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		if vulnerable {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprint(w, "panic: nil pointer dereference\n\ngoroutine 1 [running]:\n"+
				"main.handleRequest(0xc0001a4000)\n\t/srv/app/internal/backend-01.internal.corp/handlers/orders.go:42 +0x1a5\n"+
				"host: backend-01.internal.corp")
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "internal server error"}`))
		}
	})

	// vulnerable: an internal management dashboard is reachable from the
	// public interface with no authentication
	http.HandleFunc("/actuator/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if vulnerable {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "UP"}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	// vulnerable: directory listing is enabled on a proxy-served static path
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if vulnerable {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<html><body><h1>Index of /static/</h1><ul>` +
				`<li><a href="config.yml">config.yml</a></li>` +
				`<li><a href="backup.sql">backup.sql</a></li>` +
				`</ul></body></html>`))
		} else {
			http.NotFound(w, r)
		}
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
