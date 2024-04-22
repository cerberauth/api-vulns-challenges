package serve

import (
	"log"
	"net/http"
)

func RunServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
