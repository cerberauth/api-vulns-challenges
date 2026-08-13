package serve

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

const LegitKid = "legit-key-1"

// parseCertDER extracts the RSA public key from a single DER-encoded
// certificate. No chain-of-trust validation is performed against any
// pinned CA or keystore: whichever certificate is supplied is trusted.
func parseCertDER(der []byte) (*rsa.PublicKey, error) {
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, err
	}

	pub, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("certificate does not contain an RSA public key")
	}
	return pub, nil
}

// publicKeyFromX5C reads the leaf certificate out of the JWT "x5c" header
// (RFC 7515 §4.1.6) and trusts it as-is.
func publicKeyFromX5C(header map[string]interface{}) (*rsa.PublicKey, error) {
	raw, ok := header["x5c"].([]interface{})
	if !ok || len(raw) == 0 {
		return nil, fmt.Errorf("missing x5c header")
	}

	leaf, ok := raw[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid x5c entry")
	}

	der, err := base64.StdEncoding.DecodeString(leaf)
	if err != nil {
		return nil, err
	}
	return parseCertDER(der)
}

// publicKeyFromX5U fetches whatever certificate lives at the attacker
// (or legitimate) supplied "x5u" URL and trusts it as-is.
func publicKeyFromX5U(header map[string]interface{}) (*rsa.PublicKey, error) {
	x5uURL, ok := header["x5u"].(string)
	if !ok || x5uURL == "" {
		return nil, fmt.Errorf("missing x5u header")
	}

	resp, err := http.Get(x5uURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("x5u endpoint returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if block, _ := pem.Decode(body); block != nil {
		return parseCertDER(block.Bytes)
	}
	return parseCertDER(body)
}

func RunServer(port string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	legitCertBytes, err := os.ReadFile(path.Join(cwd, "keys", "legit_cert.pem"))
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(legitCertBytes)
	if block == nil {
		log.Fatal("failed to decode legit certificate PEM")
	}

	legitPublicKey, err := parseCertDER(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/.well-known/legit_cert.pem", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-pem-file")
		w.Write(legitCertBytes)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(401)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Trust whatever certificate material the caller supplies in
			// the header instead of validating it against a pinned
			// certificate or keystore.
			if pub, err := publicKeyFromX5C(token.Header); err == nil {
				return pub, nil
			}
			if pub, err := publicKeyFromX5U(token.Header); err == nil {
				return pub, nil
			}

			return legitPublicKey, nil
		})

		if err == nil && token.Valid {
			w.WriteHeader(204)
		} else {
			fmt.Println(err)
			w.WriteHeader(401)
		}
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
