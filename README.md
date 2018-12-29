# GoAuth Server

This is an implementation of an [OAuth2 Server](https://tools.ietf.org/html/rfc6749) in [Golang](https://golang.org) using [Go-kit](https://gokit.io/). It's a bit of a pet project, not intended for production (on your own head be it).

Currently it integrates with PostgreSQL, and will require PostgreSQL to be running on your local machine for the unit tests to run.

## Contributing

No contributions are being accepted at this point. It's really just a pet project to learn Golang!

## Developing, building and installing

Building/installing it:-

```bash
$ go install
```

Running it (as long as you have $GOPATH set):-

```bash
$ goauth
```

Run the unit tests:
`go test ./...`

To look at the (frankly shocking) coverage:

`go test ./... -cover`

To get some HTML for the shockingly low level of coverage:
`go test ./... -cover -coverprofile=c.out && go tool cover -html=c.out -o coverage.html && open coverage.html`

To test manually:

```bash
curl -X POST \
  http://localhost:8080/token \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -H 'Authorization: Basic Zm9vOnNlY3JldA==' \
  -H 'cache-control: no-cache' \
  -d 'grant_type=client_credentials&scope=read'
```

  Packages used:-
```bash
  go get github.com/stretchr/testify
  go get github.com/go-kit
  go get github.com/prometheus/client_golang
  go get github.com/dgrijalva/jwt-go
  go get github.com/google/uuid
  go get github.com/lib/pq
  github.com/VividCortex/gohistogram // this is needed for testing :(
```

# Integrating with PostgreSQL

Run Postgres itself:-

`./scripts/runPostgres.sh`

Connect to it using psql or similar and run contents of

`./scripts/createPostgresDb.sql`

Use this connection string:

`postgres://postgres:password@localhost/goauth`

## Endpoints

- `/token` : to retrieve either a bearer or reference token dependent on client. 200 if OK, 401 if unauthorized, 500 for anything else.
- `/validate` : used to validate a reference token, returns 200 if it's OK, or 400 if invalid.
- `/metrics`: Prometheus metrics endpoints. Full of goodness.