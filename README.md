# GoAuth Server

This is an implementation of an [OAuth2 Server](https://tools.ietf.org/html/rfc6749) in [Golang](https://golang.org) using [Go-kit](https://gokit.io/). It's a bit of a pet project, not intended for production (on your own head be it).

Currently it integrates with PostgreSQL and will require PostgreSQL to be running on your local machine for the unit tests to run (you can run tests on it using Docker, if you like).

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

  Packages used:-

```bash
  go get github.com/stretchr/testify
  go get github.com/go-kit
  go get github.com/prometheus/client_golang
  go get github.com/dgrijalva/jwt-go
  go get github.com/google/uuid
  go get github.com/lib/pq
```

# Running the Service in Docker

There is a `Dockerfile` to run the server locally with in-memory settings. There is also a `docker-compose.yml` file for running the server with a local PostgreSQL instance, which is handy for tests. All you need to do is run this in the main directory:-

`docker-compose up`

You should have th server accessible on port 8080.

# Outside-In Testing the Service in Docker

You can run the outside-in tests using:-

```bash
./test.sh
```

This will:-

- Spin up a container of PostgreSQL
- Spin up a container of GoAuth
- Spin up a container that runs nodeJS tests against GoAuth

The tests are written in nodeJS, and are pretty simple. Theoretically they could point to a different implementation of an OAuth2 service.

## Unit and Manual Testing

To run the unit tests you will need an instance of PostgreSQL to connect to. To do this, run PostgreSQL using Docker:-

`./runLocalPostgres.sh`

And use this connection string:

`postgres://postgres:password@localhost/goauth`

If you don't want to use the `./runLocalPostgres.sh` script, or you have a local instance of PostgreSQL you use for other reasons and wish to reuse it, then look in the `scripts/` folder for more information on how to set up a local instance.

### Unit Tests

You can now run the unit tests:
`go test ./...`

To look at the (frankly shocking) coverage:

`go test ./... -cover`

To get some HTML for the shockingly low level of coverage:
`go test ./... -cover -coverprofile=c.out && go tool cover -html=c.out -o coverage.html && open coverage.html`

### Manual Tests

It is often useful to do some manual testing using cURL or a UI like Postman. Below is a sample test, but in future a set of scripts will be made available so you can try all of teh functions yourself.

```bash
curl -X POST \
  http://localhost:8080/token \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -H 'Authorization: Basic Zm9vOnNlY3JldA==' \
  -H 'cache-control: no-cache' \
  -d 'grant_type=client_credentials&scope=read'
```

## Load Tests

Currently there are no load tests written. Do not use in production!

## Endpoints

An OpenApi document will be created for v1. However, in the meantime...

- `/token` : to retrieve either a bearer or reference token dependent on client. 200 if OK, 401 if unauthorized, 500 for anything else.
- `/validate` : used to validate a reference token, returns 200 if it's OK, or 400 if invalid.
- `/metrics`: Prometheus metrics endpoints. Full of goodness.