package auth

import (
	"context"
	ssov1 "sso/internal/grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface{
	Login(ctx context.Context, email string, password string, appId int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userId int64, err error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type serverApi struct{
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth){
	ssov1.RegisterAuthServer(gRPC, &serverApi{auth: auth})
}

func (s *serverApi) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error){
	if req.GetEmail() == ""{
		return nil, status.Error(codes.InvalidArgument, "message is required")
	}
	if req.GetPassword() == ""{
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.AppId == 0{
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}
	token,  err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil{
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}
func (s *serverApi) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error){
	if req.GetEmail() == ""{
		return nil, status.Error(codes.InvalidArgument, "message is required")
	}
	if req.GetPasswoed() == ""{
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userId, err :=  s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPasswoed())
	if err != nil{
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userId,
	}, nil
}



func (s *serverApi) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error){
	if req.GetUserId() == 0{
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil{
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}


