package serve

import (
	"crypto/tls"
	"log"
	"net/http"
)

func RunServer(port string, vulnerable bool) {
	cert, err := leafCertificate(vulnerable)
	if err != nil {
		log.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Write([]byte(`{"message": "ok"}`))
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		},
	}

	log.Println("Server started at port", port)
	log.Fatal(server.ListenAndServeTLS("", ""))
}
