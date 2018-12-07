See https://tools.ietf.org/html/rfc6749

TODO:

- Tighten signing of JWT
  - Validate Client secret & client Id
  - Validate Client scopes
  - Use client data to store the expiry time
- pass in the location of the .pem key files
- write script to generate .pem (openssl) - with password
- use jwt/ParseRSAPrivateKeyFromPEMWithPassword
- Add JWKS endpoint (to support Bearer tokens)
- Add verification endpoint that supports multiple certs (to enable rollover)

- create dockerfile for build

- refactor into /cmd, /pkg etc

- Add repository for clients/client credentials, scopes. Back with postgres?
- Locally cache clients in memory.
  - Have cache reload mechanism to auto pull client details. Could we have a hash in the DB to detect this?
  - Have some kind of push mechanism? OR is this simply not required if reload is efficient enough to check often(e.g. every 10 seconds).
- Add in functionality for reference tokens (poss use other cache?)

- Add OpenTracing
- Add etcd/consul for config

Running it:-

```bash
$ goauth
```

Run the unit tests:-
`go test`

To look at the (frankly shocking) coverage:-

`go test -cover`

To get some HTML foe the shockingly low level of coverage:-
`go test -cover -coverprofile=c.out && go tool cover -html=c.out -o coverage.html && open coverage.html`

To test manually:-
curl -X POST \
  http://localhost:8080/token \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -H 'Authorization: Basic bXlBd2Vzb21lQ2xpZW50OnN1cGVyc2VjcmV0cGFzc3dvcmQ=' \
  -H 'cache-control: no-cache' \
  -d 'grant_type=client_credentials&scope=read'

  Packages used:-
  go get github.com/stretchr/testify
  go get github.com/go-kit
  go get github.com/prometheus/client_golang
  go get github.com/dgrijalva/jwt-go
  go get github.com/google/uuid