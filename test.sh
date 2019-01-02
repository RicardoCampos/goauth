#!/bin/bash
set -x

cleanup () {
  # Kill all remaining processes
  docker stop $(docker-compose ps -q)
  docker-compose -p ci kill
  docker-compose -p ci rm -f

  # Return .dockerignore file back to normal
  if [ -f ./no ]; then
    mv ./no ./.dockerignore
  fi

  # Final check for remaining processes
  docker-compose ps
}

docker-compose -f docker-compose-test.yml build

# Build the test container
docker build --no-cache -f ./tests/Dockerfile-test -t goauth-test ./tests/ 

# run tests in docker container against the docker-composed services.
#docker run --rm -e TEST_ENDPOINT=http://host.docker.internal:8080 goauth-test; test_exit_code=$?
docker-compose -f docker-compose-test.yml run --rm goauth-test; test_exit_code=$?

cleanup
exit $test_exit_code