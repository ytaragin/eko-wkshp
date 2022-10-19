## Prerequisites
1. Install Go 1.18 or above (next step will verify this)
   > We can use an earlier version if 1.18 is an issue (need to download an earlier version of `sqlc` to do so)
2. Install DB tools using `make install-tools`
3. Edit the Makefile's _PostgreSQL parameters_ section to match your DB configuration. 
   > Right now it's configured to use CCS-Dev's defaults with a DB called _protection_.
4. Create a new database named `DBNAME` (e.g. SSH into the PostgreSQL container and run `createdb --username=<PG_USER> --owner=<PG_USER> <DBNAME>`)
5. Create (migrate) the database scheme using `make migrate-up`. **Note:** make sure you run the make target from an instance that can reach `(PG_HOST)`.
6. Run `go mod tidy` to download any 3rd party libraries used in code.


## Development
1. To generate a new SQL-client code, execute `sqlc generate -f ./db/sqlc.yaml`.
   > generated code should be written to `./db/repository/`.
2. execute the program using `go run ./cmd/main.go`