
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

ssh -T git@github.com
git remote set-url origin git@github.com:tonghoameo/goknab.git


reference 
https://github.com/gin-gonic/gin/blob/master/docs/doc.md#using-get-post-put-patch-delete-and-options


run test with curl

create account
curl -X POST -d '{"owner":"asdsad","currency":"USD"}' http://localhost:8888/accounts


how to use mock with golang

go mod tidy
go get github.com/golang/mock
go mod download github.com/golang/mock
mockgen -destination db/mock/store.go github.com/binbomb/goapp/simplebank/db/sqlc Store

sudo docker run --network="host" --name simplebank -p 8888:8888 -e GIN_MODE=release -e DB_URI=postgres://root:secret@127.0.0.1:5432/simple_bank?sslmode=disable simplebank:latest
mockgen -package mockdb -destination db/mock/store.go github.com/binbomb/goapp/simplebank/db/sqlc Store


dbml 
https://dbdocs.io/duyduymeo/simplebank?view=relationships

protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
proto/*.proto


https://github.com/googleapis/googleapis

cp from googleapis/google/api to proto/google/api

edit service simplebank to run grpc and http once code

add --grpc-gateway_out=pb


add swagger https://app.swaggerhub.com/


using docker compose export ports in each services

docker volume ls
check ingress-nginx

kubectl get pods -n ingress-nginx

kubectl create deployment nginx-deployment --image=nginx       
deployment.apps/nginx-deployment created
                                                                                                                                                                       
kubectl create service nodeport nginx-service --tcp=80:80
service/nginx-service created
minikube addons enable ingress
