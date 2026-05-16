package account

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/Nitish0007/go_notifier/internal/shared/api"
	accountv1 "github.com/Nitish0007/go_notifier/pkg/gen/account/v1"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type AccountGRPCServer struct {
	accountv1.UnimplementedAccountServiceServer
	svc *AccountService
}

func NewAccountServiceServer(accountService *AccountService) accountv1.AccountServiceServer {
	return &AccountGRPCServer{
		svc: accountService,
	}
}

func (s *AccountGRPCServer) Signup(ctx context.Context, req *accountv1.SignupRequest) (*accountv1.SignupResponse, error) {
	payload := &SignupRequest{}
	payload.Account.Email = req.GetEmail()
	payload.Account.Password = req.GetPassword()
	payload.Account.ConfirmPassword = req.GetConfirmPassword()
	payload.Account.FirstName = req.GetFirstName()
	payload.Account.LastName = req.GetLastName()

	validatedPayload, err := api.ValidateRequestData(payload)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	response, err := s.svc.CreateAccount(ctx, validatedPayload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create account: %v", err)
	}

	return &accountv1.SignupResponse{
		Id: response.ID,
		Email: response.Email,
		FirstName: response.FirstName,
		LastName: response.LastName,
		CreatedAt: timestamppb.New(response.CreatedAt),
		UpdatedAt: timestamppb.New(response.UpdatedAt),
	}, nil
}

func (s *AccountGRPCServer) Login(ctx context.Context, req *accountv1.LoginRequest) (*accountv1.LoginResponse, error) {
	payload := &LoginRequest{}
	payload.Login.Email = req.GetEmail()
	payload.Login.Password = req.GetPassword()

	validatedPayload, err := api.ValidateRequestData(payload)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	response, err := s.svc.Login(ctx, validatedPayload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login: %v", err)
	}
	return &accountv1.LoginResponse{
		AuthToken: response.AuthToken,
		AccountId: response.AccountID,
		Email: response.Email,
		FirstName: response.FirstName,
		LastName: response.LastName,
	}, nil
}