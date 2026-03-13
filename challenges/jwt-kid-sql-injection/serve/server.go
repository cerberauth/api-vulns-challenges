package serve

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/cerberauth/api-vulns-challenges/common"
	"github.com/golang-jwt/jwt/v5"
	_ "modernc.org/sqlite"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE keys (kid TEXT PRIMARY KEY, secret TEXT NOT NULL)`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`INSERT INTO keys (kid, secret) VALUES ('default', 'supersecretkey_stored_in_database')`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunServer(port string) {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := common.ExtractBearerToken(r)
		if !ok {
			w.WriteHeader(401)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header["kid"].(string)
			if !ok || kid == "" {
				return nil, fmt.Errorf("missing kid header")
			}

			// VULNERABILITY: SQL injection via unsanitized kid header
			query := "SELECT secret FROM keys WHERE kid = '" + kid + "'"
			var secret string
			err := db.QueryRow(query).Scan(&secret)
			if err != nil {
				return nil, fmt.Errorf("key not found: %v", err)
			}

			return []byte(secret), nil
		})

		if token != nil && token.Valid {
			w.WriteHeader(204)
		} else {
			fmt.Println(err)
			w.WriteHeader(401)
		}
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, common.SecurityHeadersMiddleware(mux)))
}
