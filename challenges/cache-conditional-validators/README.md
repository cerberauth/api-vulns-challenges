# Conditional Request / Validator Misuse

This challenge demonstrates broken cache validators: `/resource` never changes, but the vulnerable configuration issues a brand-new, random `ETag` on every response (a "flapping" validator) and ignores `If-None-Match` / `If-Modified-Since` entirely — so conditional requests never get a `304 Not Modified`, defeating revalidation and forcing a full response on every request.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /resource` — a static document that should support conditional requests

## Confirming the finding

```bash
# fetch twice and compare ETags - vulnerable: different every time; fixed: stable
curl -s -D - http://localhost:8080/resource -o /dev/null | grep ETag
curl -s -D - http://localhost:8080/resource -o /dev/null | grep ETag

# send a conditional request with the ETag just received
ETAG=$(curl -s -D - http://localhost:8080/resource -o /dev/null | grep -i etag | cut -d' ' -f2 | tr -d '\r')
curl -s -o /dev/null -w "%{http_code}\n" -H "If-None-Match: $ETAG" http://localhost:8080/resource
# expect 200 (vulnerable) vs 304 (fixed)
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: a new random ETag is issued every time and conditional requests are ignored
go run main.go serve --vulnerable=true

# fixed: a stable ETag/Last-Modified is issued and conditional requests are honored with 304
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
