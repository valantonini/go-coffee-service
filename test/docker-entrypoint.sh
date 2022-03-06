#!/bin/sh
# https://stackoverflow.com/questions/63198731/how-to-use-wait-for-it-in-docker-compose-file

# Abort on any error (including if wait-for-it fails).
set -e

# Wait for the backend to be up, if we know where it is.
/app/test/wait-for-it.sh $WAIT_FOR -t 15


# Run the main container command.
exec "$@"