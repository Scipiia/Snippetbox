package api

import "github.com/gin-gonic/gin"

type createSnippetRequest struct {
	UserID  int32  `json:"user_id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (server *Server) createSnippet(ctx *gin.Context) {

}
