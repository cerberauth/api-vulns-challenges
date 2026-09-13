# Caching Behavior & Web Cache Deception

This challenge demonstrates a shared cache that keys authenticated responses by path only (ignoring the `Authorization` header), and a web cache deception scenario where an authenticated endpoint's data is exposed under a static-looking path (`/assets/profile.js`) and gets cached and replayed to any caller.

## How to run it

```bash
go run main.go serve
```

## Endpoints

- `GET /account` — returns data scoped to the caller's `Authorization` token
- `GET /assets/profile.js` — static-looking path backed by authenticated data

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: the cache key ignores Authorization, so the first caller's response leaks to everyone; /assets/profile.js caches and serves authenticated data
go run main.go serve --vulnerable=true

# fixed: the cache key includes the Authorization token, and authenticated content is never cached or exposed under a static path
go run main.go serve --vulnerable=false
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
