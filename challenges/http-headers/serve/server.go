package serve

import (
	"log"
	"net/http"
)

func RunServer(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/wildcard-cors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
