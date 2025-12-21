package server

import (
	"context"
	"errors"

	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/card/v1"
	"github.com/aifedorov/gophkeeper/internal/server/config"
	"github.com/aifedorov/gophkeeper/internal/server/domain/auth"
	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/card"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CardServer struct {
	pb.UnimplementedCardServiceServer
	cfg     *config.Config
	logger  *zap.Logger
	authSev auth.Service
	cardSrv card.Service
}

func NewCardServer(cfg *config.Config, logger *zap.Logger, authSev auth.Service, cardSrv card.Service) *CardServer {
	return &CardServer{
		cfg:     cfg,
		logger:  logger,
		authSev: authSev,
		cardSrv: cardSrv,
	}
}

func (s *CardServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	s.logger.Debug("grpc: create card request received", zap.String("name", req.GetName()))

	userID, encryptionKey, err := s.authSev.GetUserDataFromContext(ctx)
	if err != nil {
		s.logger.Error("grpc: failed to get user ID or encryption key from token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	newCard, err := card.NewCard(
		uuid.NewString(),
		req.GetName(),
		req.GetNumber(),
		req.GetExpiredDate(),
		req.GetCardHolderName(),
		req.GetCvv(),
		req.GetNotes(),
	)
	if err != nil || newCard == nil {
		s.logger.Error("grpc: failed to create card entity", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	res, err := s.cardSrv.Create(ctx, userID, encryptionKey, *newCard)
	if errors.Is(err, card.ErrNameExists) {
		s.logger.Debug("grpc: card name already exists", zap.String("name", newCard.GetName()))
		return nil, status.Error(codes.AlreadyExists, "card name already exists")
	}
	if err != nil {
		s.logger.Error("grpc: failed to create card", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal card error")
	}

	id := res.GetID()
	s.logger.Debug("grpc: card created successfully", zap.String("id", id))

	resp := pb.CreateResponse{
		Id: &id,
	}

	return &resp, nil
}

func (s *CardServer) List(ctx context.Context, _ *pb.ListRequest) (*pb.ListResponse, error) {
	s.logger.Debug("grpc: list card request received")

	userID, encryptionKey, err := s.authSev.GetUserDataFromContext(ctx)
	if err != nil {
		s.logger.Error("grpc: failed to get user ID or encryption key from token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	cards, err := s.cardSrv.List(ctx, userID, encryptionKey)
	if err != nil {
		s.logger.Error("grpc: failed to get list of cards", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal card error")
	}

	cardList := make([]*pb.ListResponse_ListItem, len(cards))
	for i, c := range cards {
		id := c.GetID()
		name := c.GetName()
		number := c.GetNumber()
		expiredDate := c.GetExpiredDate()
		cardHolderName := c.GetCardHolderName()
		cvv := c.GetCvv()
		notes := c.GetNotes()

		cardList[i] = &pb.ListResponse_ListItem{
			Id:             &id,
			Name:           &name,
			Number:         &number,
			ExpiredDate:    &expiredDate,
			CardHolderName: &cardHolderName,
			Cvv:            &cvv,
			Notes:          &notes,
		}
	}

	return &pb.ListResponse{
		Cards: cardList,
	}, nil
}

func (s *CardServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	s.logger.Debug("grpc: update card request received")

	userID, encryptionKey, err := s.authSev.GetUserDataFromContext(ctx)
	if err != nil {
		s.logger.Error("grpc: failed to get user ID or encryption key from token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	updatedCard, err := card.NewCard(
		req.GetId(),
		req.GetName(),
		req.GetNumber(),
		req.GetExpiredDate(),
		req.GetCardHolderName(),
		req.GetCvv(),
		req.GetNotes(),
	)
	if err != nil || updatedCard == nil {
		s.logger.Error("grpc: failed to create card entity", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	res, err := s.cardSrv.Update(ctx, userID, encryptionKey, *updatedCard)
	if errors.Is(err, card.ErrNameExists) {
		s.logger.Debug("grpc: card name already exists", zap.String("name", updatedCard.GetName()))
		return nil, status.Error(codes.AlreadyExists, "card name already exists")
	}
	if errors.Is(err, card.ErrNotFound) {
		s.logger.Debug("grpc: card not found for update", zap.String("id", updatedCard.GetID()))
		return nil, status.Error(codes.NotFound, "card not found")
	}
	if err != nil {
		s.logger.Error("grpc: failed to update card", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal card error")
	}

	id := res.GetID()
	s.logger.Debug("grpc: card updated successfully", zap.String("id", id))

	success := true
	resp := pb.UpdateResponse{
		Success: &success,
	}
	return &resp, nil
}

func (s *CardServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	s.logger.Debug("grpc: delete card request received")

	userID, _, err := s.authSev.GetUserDataFromContext(ctx)
	if err != nil {
		s.logger.Error("grpc: failed to get user ID or encryption key from token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	s.logger.Debug("grpc: user_id extracted from token", zap.String("user_id", userID))

	err = s.cardSrv.Delete(ctx, userID, req.GetId())
	if errors.Is(err, card.ErrNotFound) {
		s.logger.Debug("grpc: card not found for update", zap.String("id", req.GetId()))
		return nil, status.Error(codes.NotFound, "card not found")
	}
	if err != nil {
		s.logger.Error("grpc: failed to update card", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal card error")
	}

	s.logger.Debug("grpc: card deleted successfully", zap.String("id", req.GetId()))

	success := true
	return &pb.DeleteResponse{
		Success: &success,
	}, nil
}
