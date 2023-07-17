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

	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUser)
	router.DELETE("/users/:id", server.deleteUser)
	router.PATCH("/users", server.updateUser)

	server.router = router
	return server
}

func (server *Server) Start(adress string) error {
	return server.router.Run(adress)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
