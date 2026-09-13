# Cache Poisoning via Fat GET

This challenge demonstrates a "fat GET": `/search` accepts its search term both as a `q` query parameter and, non-standardly, in the GET request body. The shared cache keys purely on the URL (query string included), since GET bodies aren't expected to influence responses, but the vulnerable origin still lets a request body override the query parameter — letting an attacker poison the cached entry for a URL with a body-controlled payload that every subsequent, body-less caller of that URL then receives.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /search?q=term` — reflects the search term, optionally overridden by the request body

## Exploiting it

```bash
# poison the cache entry for "/search?q=cats" with a body-controlled payload
curl -X GET "http://localhost:8080/search?q=cats" -d "<script>alert(1)</script>"

# any subsequent caller of the same URL, without a body, gets the poisoned response
curl "http://localhost:8080/search?q=cats"
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: a GET request body overrides the query parameter, but the cache key is unaware of the body
go run main.go serve --vulnerable=true

# fixed: the origin only ever reads the query parameter and ignores GET bodies entirely
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
