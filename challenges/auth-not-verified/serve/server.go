package serve

import (
	"log"
	"net/http"

	"github.com/cerberauth/api-vulns-challenges/common"
)

func RunServer(port string, vulnerable bool) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !vulnerable {
			// fixed: a bearer token is required and actually checked
			if _, ok := common.ExtractBearerToken(r); !ok {
				w.WriteHeader(401)
				return
			}
		}

		w.WriteHeader(204)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
