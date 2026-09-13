# HTTP Security Headers

This challenge demonstrates a reverse proxy that forwards responses without injecting common hardening headers (`X-Frame-Options`, `X-Content-Type-Options`, `Referrer-Policy`, `Permissions-Policy`) and leaks the upstream software banner via `Server` / `X-Powered-By`.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
go run main.go serve --vulnerable=true   # vulnerable: no hardening headers, Server/X-Powered-By banners disclosed
go run main.go serve --vulnerable=false  # fixed: hardening headers are injected, banners are removed
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
