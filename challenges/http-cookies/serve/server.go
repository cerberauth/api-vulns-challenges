package serve

import (
	"log"
	"net/http"
)

func RunServer(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// set unsecure cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "unsecure",
			Value: "unsecure",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
