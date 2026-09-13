# Web Cache Deception via Origin-Only Delimiters

This challenge demonstrates delimiter-based web cache deception: the shared cache treats characters like `;`, `%00`, `%0a`, or `$` as marking the end of the meaningful path, and decides to cache a request when whatever follows the delimiter looks like a static suffix (`.css`, `.js`). The origin, however, ignores the delimiter entirely and keeps serving the same dynamic, authenticated `/account` response — so the cache stores that sensitive response under a path anyone can guess.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /account[<delimiter><static-looking suffix>]` — authenticated account data, e.g. `/account;a.css`, `/account%00a.js`

## Exploiting it

```bash
# victim, authenticated, requests a delimiter-suffixed variant of /account
curl --path-as-is -H "Authorization: Bearer victim-secret-token" "http://localhost:8080/account;a.css"

# attacker, unauthenticated, requests the same path and receives the victim's cached data
curl --path-as-is "http://localhost:8080/account;a.css"
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: origin-only delimiters followed by a static-looking suffix are cached
go run main.go serve --vulnerable=true

# fixed: the authenticated response is never cached, regardless of any delimiter/suffix appended to the path
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
