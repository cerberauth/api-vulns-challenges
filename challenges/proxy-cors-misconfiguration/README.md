# CORS

This challenge demonstrates a misconfigured CORS policy on an authenticated endpoint: any `Origin` (including `null`) is reflected back combined with `Access-Control-Allow-Credentials: true`, and the preflight response is overly permissive, allowing any site to read the authenticated response in a victim's browser.

## How to run it

```bash
go run main.go serve
```

## Endpoints

- `GET /api/account` — returns account data, gated by cookies/credentials
- `OPTIONS /api/account` — CORS preflight

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: Origin is reflected with credentials allowed, null origin accepted, preflight allows any header/method
go run main.go serve --vulnerable=true

# fixed: only the trusted origin is allowed, Vary: Origin is set, preflight is scoped down
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
