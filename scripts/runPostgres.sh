#!/bin/bash
docker run --name local-postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres:11.1-alpine