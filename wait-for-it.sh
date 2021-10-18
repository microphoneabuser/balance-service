#!/bin/bash
# wait-for-postgres.sh

set -e

host="$1"
shift
cmd="$@"

until PGPASSWORD=$DB_PASSWORD psql -h "$host" -U "postgres" -c '\l'; do >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

while ! nc -z rabbitmq 5672;do >&2 echo "RabbitMQ is unavailable - sleeping"
  sleep 3
done

>&2 echo "Postgres & RabbitMQ is up - executing command"
exec $cmd