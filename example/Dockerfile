FROM previousnext/mysql:alpine

# Add the database from the compiled source.
ADD db.sql /tmp/db.sql

# Import with upstream scripts.
RUN database-import developer password local /tmp/db.sql