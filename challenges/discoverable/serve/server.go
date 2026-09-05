package serve

import (
	"log"
	"net/http"

	"github.com/cerberauth/api-vulns-challenges/common"
)

func RunServer(port string, vulnerable bool) {
	mux := http.NewServeMux()

	if vulnerable {
		// vulnerable: the OpenAPI spec (and any other file under ./static)
		// is served publicly, letting anyone enumerate the full API surface
		fs := http.FileServer(http.Dir("./static"))
		mux.Handle("/", fs)
	}

	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
