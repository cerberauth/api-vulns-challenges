package serve

import (
	"log"
	"net/http"

	"github.com/cerberauth/api-vulns-challenges/common"
)

func RunServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
