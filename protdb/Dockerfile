# splitting this stage to avoid installing cURL and migrate every time the migrations directory changes
FROM alpine 
RUN apk --no-cache add curl
# download the release of migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar -xvz
#FROM alpine-with-migrate AS migrator
# you may need to override these depending on your setup
ENV PG_USER=postgres
ENV PG_PASSWORD=mysecret
ENV PG_HOST=wkshp-postgresql
ENV PG_PORT=5432
ENV DBNAME=protection
# copy the migration files (may need change source path depending on your execution context)
COPY db/migrations /migrations


CMD /migrate -path /migrations -database "postgresql://$PG_USER:$PG_PASSWORD@$PG_HOST:$PG_PORT/$DBNAME?sslmode=disable" -verbose up
#CMD sleep 3600











# FROM golang:1.16-alpine

# WORKDIR /usr/src/app

# RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# RUN go install github.com/kyleconroy/sqlc/cmd/sqlc@latest





