package serve

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net/http"
	"time"
)

// selfSignedCert generates an ephemeral self-signed certificate so the
// challenge does not depend on any external files.
func selfSignedCert() tls.Certificate {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	tmpl := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"localhost"},
	}
	der, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &key.PublicKey, key)
	if err != nil {
		log.Fatal(err)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Fatal(err)
	}
	return cert
}

func RunServer(port string, vulnerable bool) {
	cert := selfSignedCert()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if !vulnerable {
			// fixed: HSTS is advertised with a long max-age, includeSubDomains and preload
			w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "ok"}`))
	})

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	if vulnerable {
		// vulnerable: obsolete protocol versions and weak, non-forward-secret
		// cipher suites are still accepted, and HSTS is never sent
		tlsConfig.MinVersion = tls.VersionTLS10
		tlsConfig.CipherSuites = []uint16{
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		}
	} else {
		// fixed: only modern, forward-secret protocol versions and ciphers are accepted
		tlsConfig.MinVersion = tls.VersionTLS12
		tlsConfig.CipherSuites = []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		}
	}

	server := &http.Server{
		Addr:      ":" + port,
		Handler:   handler,
		TLSConfig: tlsConfig,
		// disable automatic HTTP/2, which mandates a modern cipher suite set
		// incompatible with the weak configuration exercised by this challenge
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
	}

	log.Println("Server started at port", port)
	log.Fatal(server.ListenAndServeTLS("", ""))
}
