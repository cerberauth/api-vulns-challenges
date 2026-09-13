# Proxy-Level Template Injection

This challenge demonstrates a proxy templating directive (similar to Caddy's `templates` directive) that renders a reflected request header (`Referer`) through a server-side template engine, allowing an attacker to inject template directives (e.g. `{{env "SECRET_KEY"}}`) and read environment variables from the proxy process.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /` — renders the `Referer` header into the response body

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: curl -H 'Referer: {{env "SECRET_KEY"}}' http://localhost:8080/ leaks the SECRET_KEY env var
go run main.go serve --vulnerable=true

# fixed: the header value is HTML-escaped and never evaluated as a template
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
