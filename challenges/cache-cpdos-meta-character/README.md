# CPDoS: HTTP Meta Character (HMC)

This challenge demonstrates the "HTTP Meta Character" cache poisoning denial-of-service pattern: the shared cache/CDN forwards a percent-decoded query value containing control characters (NUL, CR, LF, ...) straight through, while the origin's stricter input validation rejects it with a `400 Bad Request`. The cache keys purely on the path, unaware of the query value that caused the error, so it caches the error and replays it to every subsequent caller of the same path — regardless of their own query value.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /?name=<value>` — greets the caller by name for normal requests

## Exploiting it

```bash
# a query value containing a NUL byte the origin rejects
curl -s -D - "http://localhost:8080/?name=bad%00name" -o /dev/null

# any subsequent caller of "/", even with a clean name, now gets the cached 400 error
curl -s -D - "http://localhost:8080/?name=alice" -o /dev/null
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: the origin's 400 error response is cached, keyed on path only, and replayed to later callers
go run main.go serve --vulnerable=true

# fixed: error responses are never cached, regardless of what triggered them
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
