package server

import (
	"context"
	"log"
	"net"

	"github.com/Gabriel-Schiestl/authgate/authpb"
	"github.com/Gabriel-Schiestl/authgate/internal/src/application/dtos"
	"github.com/Gabriel-Schiestl/authgate/internal/src/controller"
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/models"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type AuthServiceServer struct {
	authpb.UnimplementedAuthServiceServer
	controller *controller.Controller
}

func NewAuthServiceServer(lc fx.Lifecycle, controller *controller.Controller) *AuthServiceServer {
    server := &AuthServiceServer{
        controller: controller,
    }

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            lis, err := net.Listen("tcp", ":50051")
            if err != nil {
                return err
            }

            grpcServer := grpc.NewServer()
            authpb.RegisterAuthServiceServer(grpcServer, server)

            log.Printf("gRPC server listening at %v", lis.Addr())
            
            go func() {
                if err := grpcServer.Serve(lis); err != nil {
                    log.Fatalf("Failed to serve: %v", err)
                }
            }()

            return nil
        },
        OnStop: func(ctx context.Context) error {
            log.Println("Stopping gRPC server...")
            return nil
        },
    })

    return server
}

func (s *AuthServiceServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	response, err := s.controller.Login(ctx, dtos.LoginDTO{
		IdentifierType: models.IdentifierType(req.GetIdentifierType()),
		IdentifierValue: req.GetIdentifierValue(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}
	return &authpb.LoginResponse{
		Success: true,
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		UserInfo: &authpb.UserInfo{
			UserId: response.UserInfo.UserID,
			Name: response.UserInfo.Name,
			Roles:  response.UserInfo.Roles,
		},
	}, nil
}

func (s *AuthServiceServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	maxTokenAge := int(req.GetMaxTokenAgeSeconds())
	maxWrongAttempts := int(req.GetMaxWrongAttempts())

	response, err := s.controller.Register(ctx, dtos.RegisterDTO{
		IdentifierType: models.IdentifierType(req.GetIdentifierType()),
		IdentifierValue: req.GetIdentifierValue(),
		Password: req.GetPassword(),
		UserInfo: dtos.UserInfoDTO{
			Name:  req.GetUserInfo().GetName(),
			Roles: req.GetUserInfo().GetRoles(),
			UserID: req.GetUserInfo().GetUserId(),
		},
		EncryptToken:       req.GetEncryptToken(),
		MaxTokenAgeSeconds: &maxTokenAge,
		MaxWrongAttempts:   &maxWrongAttempts,
	})
	if err != nil {
		return nil, err
	}
	return &authpb.RegisterResponse{
		Success: true,
		IdentifierType: authpb.IdentifierType(response.IdentifierType),
		IdentifierValue: response.IdentifierValue,
		UserInfo: &authpb.UserInfo{
			UserId: response.UserInfo.UserID,
			Name:  response.UserInfo.Name,
			Roles: response.UserInfo.Roles,
		},
	}, nil
}

func (s *AuthServiceServer) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	response, err := s.controller.RefreshToken(ctx, dtos.RefreshTokenDTO{
		RefreshToken: req.GetRefreshToken(),
	})
	if err != nil {
		return nil, err
	}
	return &authpb.RefreshTokenResponse{
		AccessToken: response.AccessToken,
	}, nil
}

func (s *AuthServiceServer) VerifyToken(ctx context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	response, err := s.controller.VerifyToken(ctx, dtos.VerifyTokenDTO{
		AccessToken: req.GetAccessToken(),
	})
	if err != nil {
		return nil, err
	}
	return &authpb.VerifyTokenResponse{
		Success: true,
		UserInfo: &authpb.UserInfo{
			UserId: response.UserID,
			Name:  response.Name,
			Roles: response.Roles,
		},
	}, nil
}

func (s *AuthServiceServer) DeleteAuth(ctx context.Context, req *authpb.DeleteAuthRequest) (*authpb.DeleteAuthResponse, error) {
	err := s.controller.DeleteAuth(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &authpb.DeleteAuthResponse{
		Success: true,
	}, nil
}

