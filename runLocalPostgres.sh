#!/bin/bash
# Use this for local testing. Might be easier to use the supplied 
# Destoying old instances first
docker stop local-postgres && docker rm local-postgres && docker rmi local-postgres 
docker build -t local-postgres -f ./Dockerfile-postgres-test .
docker run --name local-postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -d local-postgres:latest