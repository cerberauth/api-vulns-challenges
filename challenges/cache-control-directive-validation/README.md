# Conflicting / Invalid Cache-Control Directives

This challenge demonstrates a self-contradictory `Cache-Control` header: `/report` sends `no-store, max-age=3600, public, private` all at once — directives that cannot coexist per RFC 9111 (`no-store` forbids any storage while `max-age` implies storage; `public` and `private` make opposite claims about shared-cache eligibility). Different cache implementations resolve this ambiguity inconsistently, which is itself a risk for sensitive responses.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /report` — returns a sensitive JSON report

## Confirming the finding

```bash
curl -s -D - http://localhost:8080/report -o /dev/null | grep -i cache-control
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: Cache-Control combines no-store, max-age, public and private in one contradictory header
go run main.go serve --vulnerable=true

# fixed: Cache-Control unambiguously declares no-store only
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
