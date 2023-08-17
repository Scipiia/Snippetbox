package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/scipiia/snippetbox/api"
	db "github.com/scipiia/snippetbox/db/sqlc"
	_ "github.com/scipiia/snippetbox/doc/statik"
	"github.com/scipiia/snippetbox/gapi"
	"github.com/scipiia/snippetbox/pb"
	"github.com/scipiia/snippetbox/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

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
	go runGrpcServer(config, query)
	runGatewayServer(config, query)
	//runGinServer(config, query)
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
		log.Fatal("cannot create listener", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server", err)
	}
}

func runGatewayServer(config util.Config, query db.Store) {
	server, err := gapi.NewServer(config, query)
	if err != nil {
		log.Fatal("cannot create server", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSnippetboxHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	//static swagger files
	//fs := http.FileServer(http.Dir("./doc/swagger"))
	statikFs, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik fs:", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create HTTP gateway listener", err)
	}

	log.Printf("start HTTP  gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server", err)
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
