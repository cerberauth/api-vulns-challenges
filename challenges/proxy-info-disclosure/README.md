# Information Disclosure & Exposed Management Interfaces

This challenge demonstrates several information disclosure issues commonly found behind misconfigured reverse proxies: verbose error pages leaking backend internals, an unauthenticated management/health endpoint, and directory listing on a static path.

## How to run it

```bash
go run main.go serve
```

## Endpoints

- `GET /error` — triggers a simulated backend error
- `GET /actuator/health` — internal management/health endpoint
- `GET /static/` — static file listing

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: /error leaks a stack trace and internal hostname, /actuator/health is reachable, /static/ lists its contents
go run main.go serve --vulnerable=true

# fixed: /error returns a generic message, /actuator/health is not reachable, /static/ returns 404 instead of listing
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
