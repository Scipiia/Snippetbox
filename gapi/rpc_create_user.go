package gapi

import (
	"context"
	"log"

	"github.com/lib/pq"
	db "github.com/scipiia/snippetbox/db/sqlc"
	"github.com/scipiia/snippetbox/pb"
	"github.com/scipiia/snippetbox/util"
	"github.com/scipiia/snippetbox/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		// badRequest := &errdetails.BadRequest{FieldViolations: violations}
		// statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")

		// statusDetails, err := statusInvalid.WithDetails(badRequest)
		// if err != nil {
		// 	return nil, statusInvalid.Err()
		// }

		// return nil, statusDetails.Err()

		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashedPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hsah password: %s", err)
	}

	arg := db.CreateUserParams{
		Name:           req.GetName(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.query.CreateUser(ctx, arg)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			log.Println(pqError.Code.Name())
			switch pqError.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "name already exists %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (validations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateName(req.GetName()); err != nil {
		validations = append(validations, fieldValidation("name", err))
	}

	if err := validation.ValidatePassword(req.GetPassword()); err != nil {
		validations = append(validations, fieldValidation("password", err))
	}

	if err := validation.ValidateFullName(req.GetFullName()); err != nil {
		validations = append(validations, fieldValidation("full_name", err))
	}

	if err := validation.ValidateEmail(req.GetEmail()); err != nil {
		validations = append(validations, fieldValidation("email", err))
	}

	return validations
}
