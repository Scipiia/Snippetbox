package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/scipiia/snippetbox/db/sqlc"
	"github.com/scipiia/snippetbox/pb"
	"github.com/scipiia/snippetbox/util"
	"github.com/scipiia/snippetbox/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if authPayload.Name != req.GetName() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot upsate other user's info")
	}

	violations := validateUpdateUserRequest(req)
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

	// hashedPassword, err := util.HashedPassword(req.GetPassword())
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "failed to hsah password: %s", err)
	// }

	arg := db.UpdateUserParams{
		Name: req.GetName(),
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashedPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hsah password: %s", err)
		}

		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}

		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user %s", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (validations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateName(req.GetName()); err != nil {
		validations = append(validations, fieldValidation("name", err))
	}

	if req.Password != nil {
		if err := validation.ValidatePassword(req.GetPassword()); err != nil {
			validations = append(validations, fieldValidation("password", err))
		}
	}

	if req.FullName != nil {
		if err := validation.ValidateFullName(req.GetFullName()); err != nil {
			validations = append(validations, fieldValidation("full_name", err))
		}
	}

	if req.Email != nil {
		if err := validation.ValidateEmail(req.GetEmail()); err != nil {
			validations = append(validations, fieldValidation("email", err))
		}
	}

	return validations
}
