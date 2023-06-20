package gapi

import (
	"fmt"

	db "github.com/binbomb/goapp/simplebank/db/sqlc"
	"github.com/binbomb/goapp/simplebank/pb"
	"github.com/binbomb/goapp/simplebank/token"
	"github.com/binbomb/goapp/simplebank/utils"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

// create NewServer a new http server and setup routing
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create Maker token %v", err)
	}
	// find server from folder pb ex SimpleBankServer interface and implement
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	return server, nil
}
