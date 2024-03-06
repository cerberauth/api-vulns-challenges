package serve

import (
	"log"
	"net/http"
)

func RunServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// set unsecure cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "unsecure",
			Value: "unsecure",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
