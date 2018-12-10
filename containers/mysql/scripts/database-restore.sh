#!/bin/bash

# Name:        database-restore.sh
# Description: Helper script for dropping tables and importing up a database.
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

# Connection string which will be used to perform operations on the MySQL database.
CONNECTION_STRING="mysql --host=$DATABASE_HOST --port=$DATABASE_PORT --user=$DATABASE_USER --password=$DATABASE_PASSWORD $DATABASE_NAME"

TABLES=$($CONNECTION_STRING -e 'show tables' | $AWK '{ print $1}' | $GREP -v '^Tables' )

if [ "$TABLES" == "" ]
then
	echo "Error - No tables found in $DATABASE_NAME database!"
	exit 1
fi
 
for table in $TABLES
do
	echo "Deleting $DATABASE_NAME/$table"
	$CONNECTION_STRING -e "drop table $table"
done

echo "Importing database: $FILE"
$CONNECTION_STRING < $FILE