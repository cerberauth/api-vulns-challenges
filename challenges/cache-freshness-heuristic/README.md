# Heuristic Caching of Non-Cacheable Error Codes

This challenge demonstrates over-eager heuristic caching: `/transient` always returns `503 Service Unavailable` with no explicit `Cache-Control`/`Expires` header. Per RFC 9110 §15.1, only `200, 203, 204, 206, 300, 301, 404, 405, 410, 414, 501` are heuristically cacheable by default — `503` is not among them. The vulnerable cache heuristically caches it anyway, so a single transient origin failure is replayed to every caller long after the origin has recovered.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /transient` — always returns a `503` with no caching directives

## Confirming the finding

```bash
curl -s -D - http://localhost:8080/transient -o /dev/null
curl -s -D - http://localhost:8080/transient -o /dev/null
# expect X-Cache: MISS both times (fixed) vs MISS then HIT (vulnerable)
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: the 503 response is heuristically cached and replayed indefinitely
go run main.go serve --vulnerable=true

# fixed: status codes outside RFC 9110's default-cacheable list are never heuristically cached
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
