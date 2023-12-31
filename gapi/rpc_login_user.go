package gapi

import (
	"context"
	"database/sql"

	db "github.com/scipiia/snippetbox/db/sqlc"
	"github.com/scipiia/snippetbox/pb"
	"github.com/scipiia/snippetbox/util"
	"github.com/scipiia/snippetbox/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, req.GetName())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "name not found")
			// ctx.JSON(http.StatusNotFound, errorResponse(err))
			// return
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	//access token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Name,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token ")

	}

	//refresh token
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Name,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	//metadate for userAgent and ClientIP
	mtdt := server.extractMetadata(ctx)

	//create session to db
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Name:         user.Name,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (validations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateName(req.GetName()); err != nil {
		validations = append(validations, fieldValidation("name", err))
	}

	if err := validation.ValidatePassword(req.GetPassword()); err != nil {
		validations = append(validations, fieldValidation("password", err))
	}

	return validations
}
