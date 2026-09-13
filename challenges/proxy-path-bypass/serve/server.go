package serve

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

// isBlockedPath enforces the ACL. In vulnerable mode it does a naive prefix
// check on the raw, undecoded path, which is what a proxy would do if it
// inspects the path before normalization.
func isBlockedPath(r *http.Request, vulnerable bool) bool {
	path := r.URL.Path
	if vulnerable {
		return strings.HasPrefix(path, "/admin")
	}

	// fixed: decode and normalize the path the same way the backend will
	// interpret it before making the access decision
	decoded, err := url.PathUnescape(path)
	if err != nil {
		decoded = path
	}
	cleaned := "/" + strings.Trim(strings.ToLower(decoded), "/")
	for strings.Contains(cleaned, "//") {
		cleaned = strings.ReplaceAll(cleaned, "//", "/")
	}
	return strings.HasPrefix(cleaned, "/admin")
}

func RunServer(port string, vulnerable bool) {
	// a plain http.HandlerFunc, not http.ServeMux, is used directly so that
	// raw, non-normalized paths (e.g. "/ADMIN", "//admin", encoded traversal)
	// reach the handler unchanged instead of being cleaned/redirected first,
	// the same way a naive reverse proxy would forward them.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if isBlockedPath(r, vulnerable) {
			http.Error(w, `{"error": "forbidden"}`, http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "admin panel"}`))
	})

	log.Println("Server started at port", port)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
	log.Fatal(server.ListenAndServe())
}
