package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/scipiia/snippetbox/db/sqlc"
)

type createAccountRequest struct {
	Login    string `json:"login" binding:"required"`
	Username string `json:"username" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Login:    req.Login,
		Username: req.Username,
	}

	user, err := server.query.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.query.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//отправляем пустые данные что бы тест провалился
	//user = db.User{}

	ctx.JSON(http.StatusOK, user)
}

type deleteUserRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.query.DeleteUser(ctx, req.ID)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	ctx.JSON(http.StatusNotFound, errorResponse(err))
		// 	return
		// }
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

type updateUserRequest struct {
	ID       int32  `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		ID:       req.ID,
		Username: req.Username,
	}

	user, err := server.query.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
