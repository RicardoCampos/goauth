#!/bin/bash
docker run --rm -v "$PWD":/usr/src/goauth -w /usr/src/goauth golang:1.11 go build -v