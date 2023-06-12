
designed database on dbdiagrame 
add not null becaused go string is not NULL
paste code in file export to migrations folder
check setting in sqlc.yaml
sqlx migrate add init -r

sqlx migrate run --database-url
sqlx migrate revert --database-url

using createdb and dropdb to resert if migration error

pgcli or sudo -u postgres psql

run sqlc generate

query force -- name CreateAccount :one


note select for no key update to update data
