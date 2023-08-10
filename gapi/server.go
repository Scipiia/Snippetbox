package gapi

import (
	"fmt"

	db "github.com/scipiia/snippetbox/db/sqlc"
	"github.com/scipiia/snippetbox/pb"
	"github.com/scipiia/snippetbox/token"
	"github.com/scipiia/snippetbox/util"
)

type Server struct {
	config util.Config
	query  db.Store //mock
	//token
	tokenMaker token.Maker
	pb.UnimplementedSnippetboxServer
}

// *db.Queries change on db.Store mock db
func NewServer(config util.Config, query db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKye)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		query:      query,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
