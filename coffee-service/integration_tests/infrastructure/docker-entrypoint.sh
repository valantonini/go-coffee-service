#!/bin/sh
# https://stackoverflow.com/questions/63198731/how-to-use-wait-for-it-in-docker-compose-file

# Abort on any error (including if wait-for-it fails).
set -e

# Wait for the backend to be up, if we know where it is.
#if [ -n "$CUSTOMERS_HOST" ]; then
  /app/integration_tests/infrastructure/wait-for-it.sh "api:8080" -t 5
#fi

# Run the main container command.
exec "$@"