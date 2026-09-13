package serve

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
)

const (
	fixedBody     = "hello, world"
	fixedETag     = `"a1b2c3d4"`
	lastModified  = "Wed, 21 Oct 2020 07:28:00 GMT"
	contentTypeVl = "text/plain"
)

func randomETag() string {
	b := make([]byte, 4)
	rand.Read(b)
	return `"` + hex.EncodeToString(b) + `"`
}

func RunServer(port string, vulnerable bool) {
	// /resource is a static, unchanging document. A well-behaved origin
	// should hand out a stable validator (ETag/Last-Modified) and honor
	// conditional requests with 304 Not Modified. The vulnerable
	// configuration issues a brand-new, random ETag on every response
	// (a "flapping" validator) and never checks If-None-Match /
	// If-Modified-Since at all, defeating revalidation entirely - every
	// conditional request still gets a full 200 response.
	http.HandleFunc("/resource", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentTypeVl)

		if vulnerable {
			w.Header().Set("ETag", randomETag())
			w.Write([]byte(fixedBody))
			return
		}

		w.Header().Set("ETag", fixedETag)
		w.Header().Set("Last-Modified", lastModified)

		if inm := r.Header.Get("If-None-Match"); inm != "" && inm == fixedETag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		if ims := r.Header.Get("If-Modified-Since"); ims != "" && ims == lastModified {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		w.Write([]byte(fixedBody))
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
