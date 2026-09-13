package serve

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	textTemplate "text/template"
)

// funcs exposes a couple of primitives similar to what a proxy's templating
// directive (e.g. Caddy's `templates`) would offer, including reading files
// and environment variables.
var funcs = textTemplate.FuncMap{
	"env": os.Getenv,
}

func RunServer(port string, vulnerable bool) {
	if os.Getenv("SECRET_KEY") == "" {
		os.Setenv("SECRET_KEY", "supersecret-internal-key")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		referer := r.Header.Get("Referer")
		w.Header().Set("Content-Type", "text/html")

		if vulnerable {
			// vulnerable: the Referer header is rendered through the proxy's
			// server-side templating engine, so an attacker-controlled value
			// like "{{env \"SECRET_KEY\"}}" gets evaluated instead of
			// displayed as plain text
			tmpl, err := textTemplate.New("page").Funcs(funcs).Parse(
				`<html><body>Referred by: ` + referer + `</body></html>`,
			)
			if err != nil {
				http.Error(w, "template error", http.StatusInternalServerError)
				return
			}
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, nil); err != nil {
				http.Error(w, "template error", http.StatusInternalServerError)
				return
			}
			w.Write(buf.Bytes())
			return
		}

		// fixed: the value is HTML-escaped and never parsed as a template,
		// so template directives are displayed as inert text
		tmpl := template.Must(template.New("page").Parse(
			`<html><body>Referred by: {{.}}</body></html>`,
		))
		var buf bytes.Buffer
		tmpl.Execute(&buf, referer)
		w.Write(buf.Bytes())
	})

	log.Println("Server started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
