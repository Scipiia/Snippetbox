package gapi

import (
	"context"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	db "github.com/scipiia/snippetbox/db/sqlc"
	"github.com/scipiia/snippetbox/pb"
	"github.com/scipiia/snippetbox/util"
	"github.com/scipiia/snippetbox/validation"
	"github.com/scipiia/snippetbox/worker"
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

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Name:           req.GetName(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Name: user.Name,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskDistributer.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := server.store.CreateUserTx(ctx, arg)
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
		User: convertUser(txResult.User),
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
