FROM previousnext/mysql:latest

# Add the database from the compiled source.
ADD db.sql /tmp/db.sql

# Import with upstream scripts.
RUN database-import drupal drupal local /tmp/db.sql
