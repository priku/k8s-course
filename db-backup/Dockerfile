# The same Postgres image as the database, so pg_dump has the same major version as the server
FROM postgres:18.6-alpine3.24

# curl for the token exchange and the upload, jq to read the token from the response
RUN apk add --no-cache curl jq

COPY backup.sh /usr/local/bin/backup.sh

# The postgres user of the image (uid 70) is enough: the dump is written to /tmp
USER postgres

ENTRYPOINT ["/usr/local/bin/backup.sh"]
