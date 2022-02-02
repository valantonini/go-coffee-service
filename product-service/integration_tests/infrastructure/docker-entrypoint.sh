#!/bin/sh
# https://stackoverflow.com/questions/63198731/how-to-use-wait-for-it-in-docker-compose-file

# Abort on any error (including if wait-for-it fails).
set -e

# Wait for the backend to be up, if we know where it is.

/app/integration_tests/infrastructure/wait-for-it.sh $WAIT_FOR_NATS -t 5
/app/integration_tests/infrastructure/wait-for-it.sh $WAIT_FOR -t 5


# Run the main container command.
exec "$@"