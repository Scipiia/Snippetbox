package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/scipiia/snippetbox/api"
	db "github.com/scipiia/snippetbox/db/sqlc"
	"github.com/scipiia/snippetbox/gapi"
	"github.com/scipiia/snippetbox/pb"
	"github.com/scipiia/snippetbox/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

// const (
// 	dbDriver     = "postgres"
// 	dbSource     = "postgresql://root:secret@localhost:5432/snippetbox?sslmode=disable"
// 	serverAdress = "0.0.0.0:8080"
// )

func main() {
	config, err := util.LiadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	query := db.New(conn)
	runGrpcServer(config, query)
}

func runGrpcServer(config util.Config, query db.Store) {
	server, err := gapi.NewServer(config, query)
	if err != nil {
		log.Fatal("cannot create server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSnippetboxServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}

func runGinServer(config util.Config, query db.Store) {
	server, err := api.NewServer(config, query)
	if err != nil {
		log.Fatal("cannot create server", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
