package server

import (
	"context"
	"errors"

	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/auth/v1"
	"github.com/aifedorov/gophkeeper/internal/server/config"
	userdomain "github.com/aifedorov/gophkeeper/internal/server/domain/user"
	"github.com/aifedorov/gophkeeper/internal/server/infrastructure/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	errMsgInvalidCredentials = "invalid login or password"
	errMsgLoginExists        = "login already exists"
	errMsgInternalError      = "internal server error"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	cfg     *config.Config
	logger  *zap.Logger
	userSrv userdomain.Service
	jwtSrv  jwt.Service
}

// NewAuthServer creates a new instance of AuthServer with the provided dependencies.
// It initializes the gRPC authentication server that handles user registration and login.
func NewAuthServer(cfg *config.Config, logger *zap.Logger, userSrv userdomain.Service, jwtSrv jwt.Service) *AuthServer {
	return &AuthServer{
		cfg:     cfg,
		logger:  logger,
		userSrv: userSrv,
		jwtSrv:  jwtSrv,
	}
}

// Register handles user registration requests.
// It validates credentials, creates a new user account, and returns a JWT access token.
// Returns an error with gRPC status code AlreadyExists if the login already exists,
// InvalidArgument if credentials are invalid, or Internal if an unexpected error occurs.
func (a *AuthServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := a.validateCredentials(in.GetLogin(), in.GetPassword()); err != nil {
		return nil, err
	}

	a.logger.Debug("register user", zap.String("login", in.GetLogin()))

	user, err := a.userSrv.Register(in.GetLogin(), in.GetPassword())
	if errors.Is(err, userdomain.ErrLoginExists) {
		a.logger.Info("registration failed", zap.String("reason", errMsgLoginExists))
		return nil, status.Error(codes.AlreadyExists, errMsgLoginExists)
	}
	if err != nil {
		a.logger.Error("failed to register user", zap.Error(err))
		return nil, status.Error(codes.Internal, errMsgInternalError)
	}

	a.logger.Debug("user registered", zap.String("user_id", user.GetUserID()))

	userID := user.GetUserID()
	token, err := a.issueTokenAndLog(userID, "register")
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		UserId:      &userID,
		AccessToken: &token,
	}, nil
}

// Login handles user authentication requests.
// It validates credentials, authenticates the user, and returns a JWT access token.
// Returns an error with gRPC status code Unauthenticated if credentials are invalid,
// InvalidArgument if credentials are empty, or Internal if an unexpected error occurs.
func (a *AuthServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := a.validateCredentials(in.GetLogin(), in.GetPassword()); err != nil {
		return nil, err
	}

	a.logger.Debug("login user", zap.String("login", in.GetLogin()))
	user, err := a.userSrv.Login(in.GetLogin(), in.GetPassword())
	if errors.Is(err, userdomain.ErrUserNotFound) {
		a.logger.Info("login failed", zap.String("reason", errMsgInvalidCredentials))
		return nil, status.Error(codes.Unauthenticated, errMsgInvalidCredentials)
	}
	if err != nil {
		a.logger.Error("failed to login user", zap.Error(err))
		return nil, status.Error(codes.Internal, errMsgInternalError)
	}

	a.logger.Debug("user authenticated", zap.String("user_id", user.GetUserID()))

	userID := user.GetUserID()
	token, err := a.issueTokenAndLog(userID, "login")
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		UserId:      &userID,
		AccessToken: &token,
	}, nil
}

func (a *AuthServer) validateCredentials(login, password string) error {
	a.logger.Debug("validating credentials", zap.String("login", login))
	if login == "" || password == "" {
		a.logger.Info("invalid credentials", zap.String("reason", "empty login or password"))
		return status.Error(codes.InvalidArgument, errMsgInvalidCredentials)
	}
	return nil
}

func (a *AuthServer) issueTokenAndLog(userID, operation string) (string, error) {
	if userID == "" {
		a.logger.Error("user_id is empty", zap.String("operation", operation))
		return "", status.Error(codes.Internal, errMsgInternalError)
	}

	a.logger.Debug("issuing token", zap.String("user_id", userID), zap.String("operation", operation))

	token, err := a.jwtSrv.IssueToken(userID)
	if err != nil {
		a.logger.Error("failed to issue token", zap.Error(err), zap.String("operation", operation))
		return "", status.Error(codes.Internal, errMsgInternalError)
	}

	a.logger.Debug("token issued", zap.String("user_id", userID))
	return token, nil
}
