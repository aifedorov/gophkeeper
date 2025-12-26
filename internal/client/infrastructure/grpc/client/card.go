package client

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/card"
	pb "github.com/aifedorov/gophkeeper/internal/server/api/grpc/gen/card/v1"
	"google.golang.org/grpc"
)

type cardClient struct {
	client pb.CardServiceClient
}

func NewCardClient(conn *grpc.ClientConn) card.Client {
	return &cardClient{
		client: pb.NewCardServiceClient(conn),
	}
}

func (c *cardClient) Create(ctx context.Context, card card.Card) (id string, version int64, err error) {
	request := pb.CreateRequest{
		Name:           &card.Name,
		Number:         &card.Number,
		ExpiredDate:    &card.ExpiredDate,
		CardHolderName: &card.CardHolderName,
		Cvv:            &card.Cvv,
		Notes:          &card.Notes,
	}
	res, err := c.client.Create(ctx, &request)
	if err != nil {
		return "", 0, handleGRPCError(err)
	}

	// Note: Server doesn't return version yet, defaulting to 1
	// This will be updated when server proto is updated
	return res.GetId(), 1, nil
}

func (c *cardClient) Update(ctx context.Context, id string, card card.Card) (version int64, err error) {
	request := pb.UpdateRequest{
		Id:             &id,
		Name:           &card.Name,
		Number:         &card.Number,
		ExpiredDate:    &card.ExpiredDate,
		CardHolderName: &card.CardHolderName,
		Cvv:            &card.Cvv,
		Notes:          &card.Notes,
	}
	response, err := c.client.Update(ctx, &request)
	if err != nil {
		return 0, handleGRPCError(err)
	}

	// Note: Server doesn't return version yet, incrementing cached version
	// This will be updated when server proto is updated
	if !response.GetSuccess() {
		return 0, fmt.Errorf("client: update operation failed")
	}
	// Return incremented version for now (will be replaced with server response when proto is updated)
	return card.Version + 1, nil
}

func (c *cardClient) Delete(ctx context.Context, id string) error {
	request := pb.DeleteRequest{
		Id: &id,
	}

	resp, err := c.client.Delete(ctx, &request)
	if err != nil {
		return handleGRPCError(err)
	}
	if !resp.GetSuccess() {
		return fmt.Errorf("client: delete operation failed")
	}
	return nil
}

func (c *cardClient) List(ctx context.Context) ([]card.Card, error) {
	request := pb.ListRequest{}
	response, err := c.client.List(ctx, &request)
	if err != nil {
		return []card.Card{}, handleGRPCError(err)
	}
	cards := make([]card.Card, len(response.Cards))
	for i, c := range response.Cards {
		cards[i] = card.Card{
			ID:             c.GetId(),
			Name:           c.GetName(),
			Number:         c.GetNumber(),
			ExpiredDate:    c.GetExpiredDate(),
			CardHolderName: c.GetCardHolderName(),
			Cvv:            c.GetCvv(),
			Notes:          c.GetNotes(),
			Version:        1, // Default version until server returns it
		}
	}
	return cards, nil
}
