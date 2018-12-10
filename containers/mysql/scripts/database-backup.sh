#!/bin/bash

# Name:        database-backup.sh
# Description: Helper script for backing up a database.
# Author:      Nick Schuch

FILE=$1

if [ "$DATABASE_HOST" == "" ]; then
  echo "Not found: DATABASE_HOST"
  exit 1
fi

if [ "$DATABASE_PORT" == "" ]; then
  echo "Not found: DATABASE_PORT"
  exit 1
fi

if [ "$DATABASE_USER" == "" ]; then
  echo "Not found: DATABASE_USER"
  exit 1
fi

if [ "$DATABASE_PASSWORD" == "" ]; then
  echo "Not found: DATABASE_PASSWORD"
  exit 1
fi

if [ "$DATABASE_NAME" == "" ]; then
  echo "Not found: DATABASE_NAME"
  exit 1
fi

if [ "$FILE" == "" ]; then
  echo "Not found: FILE"
  exit 1
fi

mysqldump --single-transaction \
          --host=$DATABASE_HOST \
          --port=$DATABASE_PORT \
          --user=$DATABASE_USER \
          --password=$DATABASE_PASSWORD $DATABASE_NAME > $FILE