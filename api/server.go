package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/scipiia/snippetbox/db/sqlc"
	"github.com/scipiia/snippetbox/token"
	"github.com/scipiia/snippetbox/util"
)

type Server struct {
	config util.Config
	query  db.Store //mock
	//token
	tokenMaker token.Maker
	router     *gin.Engine
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
	// router := gin.Default()

	// //user
	// router.POST("/users", server.createUser)
	// router.POST("/users/login", server.loginUser)

	// //account
	// router.POST("/accounts", server.createAccount)
	// router.GET("/accounts/:id", server.getAccount)
	// router.DELETE("/accounts/:id", server.deleteAccount)
	// router.PATCH("/accounts", server.updateAccount)

	// //snippets
	// router.POST("/accounts/snippet", server.createSnippet)
	// router.GET("/accounts/snippet/:id", server.getSnippet)
	// router.GET("/accounts/snippet", server.listSnippets)
	// router.DELETE("/accounts/snippet/:id", server.deleteSnippet)

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	//user
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	//account
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	router.PATCH("/accounts", server.updateAccount)

	//snippets
	router.POST("/accounts/snippet", server.createSnippet)
	router.GET("/accounts/snippet/:id", server.getSnippet)
	router.GET("/accounts/snippet", server.listSnippets)
	router.DELETE("/accounts/snippet/:id", server.deleteSnippet)

	server.router = router
}

func (server *Server) Start(adress string) error {
	return server.router.Run(adress)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
