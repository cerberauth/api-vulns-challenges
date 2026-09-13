# Cache-Based Sensitive Information Disclosure

This challenge demonstrates a broad cache-based information disclosure finding: `/dashboard` issues a fresh session cookie and CSRF token per request and returns personalized JSON, but never sends `Cache-Control: private, no-store`. A shared cache in front of it stores the first caller's `Set-Cookie`, CSRF token, and personal data, and replays all of it — session included — to every subsequent caller of the same URL.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /dashboard` — returns a personalized JSON payload with a session cookie and CSRF token

## Exploiting it

```bash
# victim's session cookie and CSRF token get cached
curl -s -D - http://localhost:8080/dashboard -o /dev/null

# attacker requests the same URL and receives the victim's session cookie and CSRF token
curl -s -D - http://localhost:8080/dashboard -o /dev/null
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: no Cache-Control directive is sent, so the shared cache stores and replays session data
go run main.go serve --vulnerable=true

# fixed: Cache-Control: private, no-store is sent and the response is never cached
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
