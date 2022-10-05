#!/bin/sh

set -o errexit
set -o nounset

/app/build migration migrate
exec /app/build run
