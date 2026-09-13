package serve

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// baselinePatterns is a naive, single-pass WAF rule set. It is intentionally
// easy to bypass with encoding tricks to demonstrate an inconsistent filter.
var baselinePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)('|--|;|\bunion\b|\bselect\b|\bor\b\s+1=1)`), // SQLi
	regexp.MustCompile(`(?i)<script`),                                    // XSS
	regexp.MustCompile(`\.\./`),                                          // path traversal
}

func isMalicious(payload string) bool {
	for _, p := range baselinePatterns {
		if p.MatchString(payload) {
			return true
		}
	}
	return false
}

func RunServer(port string, vulnerable bool) {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		w.Header().Set("Content-Type", "application/json")

		if vulnerable {
			// vulnerable: no filtering at all, every payload is passed
			// through to the backend unmodified
			io.WriteString(w, `{"query": "`+q+`"}`)
			return
		}

		// fixed: the payload is normalized (URL-decoded, case-folded) before
		// being checked, so encoding-based bypass attempts are also caught
		decoded, err := url.QueryUnescape(q)
		if err != nil {
			decoded = q
		}
		normalized := strings.ToLower(decoded)
		if isMalicious(normalized) {
			http.Error(w, `{"error": "request blocked by WAF"}`, http.StatusForbidden)
			return
		}
		io.WriteString(w, `{"query": "`+q+`"}`)
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
