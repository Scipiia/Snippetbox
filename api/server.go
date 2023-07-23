package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/scipiia/snippetbox/db/sqlc"
)

type Server struct {
	query  db.Store //mock
	router *gin.Engine
}

// *db.Queries change on db.Store mock db
func NewServer(query db.Store) *Server {
	server := &Server{query: query}
	router := gin.Default()

	//user
	router.POST("/users", server.createUser)

	//account
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	router.PATCH("/accounts", server.updateAccount)

	//snippets
	router.POST("/accounts/snippet", server.createSnippet)
	router.GET("/accounts/snippet/:id", server.getSnippet)
	router.GET("/accounts/snippet", server.listSnippets)
	router.DELETE("/accounts/snippet/:id", server.deleteSnippet)

	server.router = router
	return server
}

func (server *Server) Start(adress string) error {
	return server.router.Run(adress)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
