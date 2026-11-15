package serve

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

var caSeed = sha256.Sum256([]byte("api-vulns-challenges/tls-certificate-misconfiguration"))

func CA() (*x509.Certificate, ed25519.PrivateKey, error) {
	key := ed25519.NewKeyFromSeed(caSeed[:])
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "API Vulns Challenges Test CA"},
		NotBefore:             time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		NotAfter:              time.Date(2045, 1, 1, 0, 0, 0, 0, time.UTC),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, key.Public(), key)
	if err != nil {
		return nil, nil, err
	}
	cert, err := x509.ParseCertificate(der)
	return cert, key, err
}

func CAPEM() ([]byte, error) {
	cert, _, err := CA()
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}), nil
}

func leafCertificate(vulnerable bool) (tls.Certificate, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	now := time.Now()
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(now.UnixNano()),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	var parent *x509.Certificate
	var signer interface{}
	if vulnerable {
		tmpl.Subject = pkix.Name{CommonName: "internal.example.com"}
		tmpl.DNSNames = []string{"internal.example.com"}
		tmpl.NotBefore = now.AddDate(-2, 0, 0)
		tmpl.NotAfter = now.AddDate(-1, 0, 0)
		parent, signer = tmpl, key
	} else {
		ca, caKey, err := CA()
		if err != nil {
			return tls.Certificate{}, err
		}
		tmpl.Subject = pkix.Name{CommonName: "localhost"}
		tmpl.DNSNames = []string{"localhost"}
		tmpl.IPAddresses = []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback}
		tmpl.NotBefore = now.Add(-time.Hour)
		tmpl.NotAfter = now.AddDate(0, 0, 90)
		parent, signer = ca, caKey
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, parent, &key.PublicKey, signer)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}, nil
}
