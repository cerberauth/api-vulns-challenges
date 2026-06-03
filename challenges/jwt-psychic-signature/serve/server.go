package serve

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
)

// vulnerableECDSAVerify mimics CVE-2022-21449: missing check that r,s ∈ [1,N-1].
// When s=0, the modular inverse is undefined; Java's implementation collapsed
// this to a zero point whose x-coordinate equaled r=0, making verification pass.
func vulnerableECDSAVerify(pub *ecdsa.PublicKey, hash []byte, r, s *big.Int) bool {
	if s.Sign() == 0 {
		return r.Sign() == 0
	}
	return ecdsa.Verify(pub, hash, r, s)
}

func readPublicKey() (*ecdsa.PublicKey, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := os.ReadFile(path.Join(cwd, "keys", "public_key.pem"))
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseECPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func RunServer(port string) {
	publicKey, err := readPublicKey()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(401)
			return
		}

		parts := strings.Split(tokenString, ".")
		if len(parts) != 3 {
			w.WriteHeader(401)
			return
		}

		headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
		if err != nil {
			w.WriteHeader(401)
			return
		}
		var header map[string]interface{}
		if err := json.Unmarshal(headerBytes, &header); err != nil {
			w.WriteHeader(401)
			return
		}
		if alg, _ := header["alg"].(string); alg != "ES256" {
			w.WriteHeader(401)
			return
		}

		sigBytes, err := base64.RawURLEncoding.DecodeString(parts[2])
		if err != nil || len(sigBytes) != 64 {
			w.WriteHeader(401)
			return
		}

		r2 := new(big.Int).SetBytes(sigBytes[:32])
		s2 := new(big.Int).SetBytes(sigBytes[32:])

		hash := sha256.Sum256([]byte(parts[0] + "." + parts[1]))

		if !vulnerableECDSAVerify(publicKey, hash[:], r2, s2) {
			w.WriteHeader(401)
			return
		}

		payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			w.WriteHeader(401)
			return
		}
		var claims map[string]interface{}
		if err := json.Unmarshal(payloadBytes, &claims); err != nil {
			w.WriteHeader(401)
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok || time.Now().Unix() > int64(exp) {
			fmt.Println("token expired")
			w.WriteHeader(401)
			return
		}

		w.WriteHeader(204)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
