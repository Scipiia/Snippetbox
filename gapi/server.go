package gapi

import (
	"fmt"

	db "github.com/scipiia/snippetbox/db/sqlc"
	"github.com/scipiia/snippetbox/pb"
	"github.com/scipiia/snippetbox/token"
	"github.com/scipiia/snippetbox/util"
	"github.com/scipiia/snippetbox/worker"
)

type Server struct {
	config util.Config
	store  db.Store
	//token
	tokenMaker token.Maker
	pb.UnimplementedSnippetboxServer
	taskDistributer worker.TaskDistributor
}

// *db.Queries change on db.Store mock db
func NewServer(config util.Config, store db.Store, taskDistributer worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKye)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributer: taskDistributer,
	}

	return server, nil
}
