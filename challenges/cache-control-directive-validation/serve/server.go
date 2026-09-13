package serve

import (
	"log"
	"net/http"
)

func RunServer(port string, vulnerable bool) {
	// /report returns a sensitive report. The vulnerable configuration
	// emits a self-contradictory Cache-Control header: "no-store" (never
	// persist this response) combined with "max-age=3600" (persist it for
	// an hour) and "public, private" (shareable and per-user at once).
	// RFC 9111 says a directive parser must be able to flag these as
	// invalid/conflicting, and different cache implementations resolve the
	// ambiguity inconsistently, which is itself a risk.
	http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if vulnerable {
			w.Header().Set("Cache-Control", "no-store, max-age=3600, public, private")
		} else {
			w.Header().Set("Cache-Control", "no-store")
		}
		w.Write([]byte(`{"report": "quarterly figures"}`))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
