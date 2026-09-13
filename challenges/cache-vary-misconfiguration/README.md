# Missing Vary Header Cache Poisoning

This challenge demonstrates a missing/incorrect `Vary` header: `/home` renders content that depends on the caller's session cookie, but the response omits `Vary: Cookie` and the shared cache keys purely on the path. The first cookie-bearing caller's personalized greeting is cached and served to every subsequent visitor, regardless of their own session.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /home` — greets the caller by name if a `session=<name>` cookie is present

## Exploiting it

```bash
# a logged-in user's personalized response gets cached under the shared key
curl -s -H "Cookie: session=alice" http://localhost:8080/home

# any subsequent visitor, even without a session, sees alice's greeting
curl -s http://localhost:8080/home
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: no Vary header is sent, and the cache key ignores the Cookie header entirely
go run main.go serve --vulnerable=true

# fixed: Vary: Cookie is sent and the cache key includes the Cookie header
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
