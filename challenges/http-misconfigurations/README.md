# HTTP Misconfigurations

This challenge demonstrates various HTTP misconfigurations, including CORS, CSP, and insecure cookies.

## How to run it

```bash
go run main.go
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: method override headers bypass method checks, CORS is wide open, cookies miss Secure/HttpOnly/SameSite/expiration, CSP allows framing
go run main.go serve --vulnerable=true

# fixed: each endpoint returns its hardened counterpart
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
