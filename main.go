package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/binbomb/goapp/simplebank/api"
	db "github.com/binbomb/goapp/simplebank/db/sqlc"
	"github.com/binbomb/goapp/simplebank/email"
	"github.com/binbomb/goapp/simplebank/gapi"
	"github.com/binbomb/goapp/simplebank/pb"
	"github.com/binbomb/goapp/simplebank/utils"
	"github.com/binbomb/goapp/simplebank/worker"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := utils.LoadConfig(".") // app.env
	if err != nil {
		log.Fatal().Msg("cannot load config to file: ")
	}

	if config.Environment == "developement" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal().Msg("cannot connect to db: ")
	}
	// run migrate sql
	runMigrationDB(config.MigrationURL, config.DBSource)
	//runDownDB(config.MigrationURL, config.DBSource)
	//return
	//runMigrationDB(config.MigrationURL, config.DBSource)
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	store := db.NewStore(conn)
	go runTaskProcessor(config, redisOpt, store)
	go runGatewayServer(config, store, taskDistributor)
	go runGinServer(config, store)

	runGrpcServer(config, store, taskDistributor)
}
func runGrpcServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create server grpc")
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot listen grpc")
	}

	log.Info().Msgf("start listening grpc server %s\n ", lis.Addr().String())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal().Msg("can start server")
	}

}
func runGatewayServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create server grpc")
	}
	// json fileds
	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOptions)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal().Msg("cannot register handler server ")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	lis, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listen MUX server")
	}

	log.Info().Msgf("start listening grpc server %s \n", lis.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(lis, handler)
	if err != nil {
		log.Fatal().Msg("can start server Mux")
	}

}
func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create server ")
	}
	server.Start(config.HTTPGinServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start server ")
	}
}
func runTaskProcessor(config utils.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	// create server Process
	mailer := email.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("start task process")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task process")
	}

}
func runMigrationDB(migrationURL string, dbSource string) {
	migration, err := migrate.New(
		migrationURL,
		dbSource,
	)
	if err != nil {
		log.Fatal().Msg("cannot create migration")
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msg("cannot migrate up ")
	}
	log.Info().Msg("migrate up database success")
}

func runDownDB(migrationURL string, dbSource string) {
	migration, err := migrate.New(
		migrationURL,
		dbSource,
	)
	if err != nil {
		log.Fatal().Msg("cannot create migration")
	}

	if err := migration.Down(); err != nil {
		fmt.Println("error ", err)
		log.Fatal().Msg("cannot migrate down ")
	}
	log.Info().Msg("migrate down database success")

}
