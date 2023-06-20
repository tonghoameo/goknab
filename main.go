package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/binbomb/goapp/simplebank/api"
	db "github.com/binbomb/goapp/simplebank/db/sqlc"
	"github.com/binbomb/goapp/simplebank/gapi"
	"github.com/binbomb/goapp/simplebank/pb"
	"github.com/binbomb/goapp/simplebank/utils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := utils.LoadConfig(".") // app.env
	if err != nil {
		log.Fatal("cannot load config to file: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}
func runGrpcServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server grpc", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot listen grpc", err)
	}

	fmt.Printf("start listening grpc server %s\n ", lis.Addr().String())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("can start server", err)
	}

}
func runGatewayServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server grpc", err)
	}

	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal("cannot register handler server ")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	lis, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listen MUX server", err)
	}

	fmt.Printf("start listening grpc server %s \n", lis.Addr().String())
	err = http.Serve(lis, mux)
	if err != nil {
		log.Fatal("can start server Mux", err)
	}

}
func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server ", err)
	}
	server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server ", err)
	}
}
