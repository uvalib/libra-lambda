# libra-db

libra-db uses golang style migrations to manage the db for libra.

## Usage
Required Environment:

- export DB_HOST=
- export DB_NAME=libra
- export DB_USER=libra
- export DB_PORT=5432
- export DB_PASSWORD=

Make will download migrate from <https://github.com/golang-migrate/migrate>

``` bash
# commands starting with migrate- are passed through to the migrate tool
make migrate-up
make migrate-down
make migrate-version

# Stubs out new migration .sql files in migrations/
# This does not actually run the migration
make new NAME=newMigrationName

```
