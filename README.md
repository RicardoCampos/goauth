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

# Running in Docker

There is a `Dockerfile` to run the server locally with in-memory settings. There is also a `docker-compose.yml` file for running the server with a local PostgreSQL instance, which is handy for tests. All you need to do is run this in the main directory:-

`docker-compose up`

You should have th server accessible on port 8080.

# Integrating with PostgreSQL Locally

Run PostgreSQL using Docker:-

`./runLocalPostgres.sh`

Use this connection string:

`postgres://postgres:password@localhost/goauth`

If you don't want to usethat script (it deletes all data each time), then look in the `scripts/` folder for more information on how to set up a local instance.


## Endpoints

- `/token` : to retrieve either a bearer or reference token dependent on client. 200 if OK, 401 if unauthorized, 500 for anything else.
- `/validate` : used to validate a reference token, returns 200 if it's OK, or 400 if invalid.
- `/metrics`: Prometheus metrics endpoints. Full of goodness.