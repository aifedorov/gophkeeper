package card

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/server/domain/secret/card/interfaces"
)

func toDomainCard(crypto interfaces.CryptoService, key []byte, card interfaces.RepositoryCard) (Card, error) {
	number, err := crypto.Decrypt(card.EncryptedNumber, key)
	if err != nil {
		return Card{}, fmt.Errorf("failed to decrypt number: %w", err)
	}
	expiredDate, err := crypto.Decrypt(card.EncryptedExpiredDate, key)
	if err != nil {
		return Card{}, fmt.Errorf("failed to decrypt expired date: %w", err)
	}
	cardHolderName, err := crypto.Decrypt(card.EncryptedCardHolderName, key)
	if err != nil {
		return Card{}, fmt.Errorf("failed to decrypt card holder name: %w", err)
	}
	cvv, err := crypto.Decrypt(card.EncryptedCvv, key)
	if err != nil {
		return Card{}, fmt.Errorf("failed to decrypt cvv: %w", err)
	}
	notes, err := crypto.Decrypt(card.EncryptedNotes, key)
	if err != nil {
		return Card{}, fmt.Errorf("failed to decrypt notes: %w", err)
	}
	if card.Version < 1 {
		return Card{}, fmt.Errorf("invalid card version: %d", card.Version)
	}

	return Card{
		id:             card.ID,
		userID:         card.UserID,
		name:           card.Name,
		number:         number,
		expiredDate:    expiredDate,
		cardHolderName: cardHolderName,
		cvv:            cvv,
		notes:          notes,
		version:        card.Version,
	}, nil
}

func toRepositoryCard(crypto interfaces.CryptoService, key []byte, card Card) (interfaces.RepositoryCard, error) {
	encryptNumber, err := crypto.Encrypt(card.number, key)
	if err != nil {
		return interfaces.RepositoryCard{}, fmt.Errorf("failed to encrypt number: %w", err)
	}

	encryptExpiredDate, err := crypto.Encrypt(card.expiredDate, key)
	if err != nil {
		return interfaces.RepositoryCard{}, fmt.Errorf("failed to encrypt expired date: %w", err)
	}

	encryptCardHolderName, err := crypto.Encrypt(card.cardHolderName, key)
	if err != nil {
		return interfaces.RepositoryCard{}, fmt.Errorf("failed to encrypt card holder name: %w", err)
	}

	encryptCvv, err := crypto.Encrypt(card.cvv, key)
	if err != nil {
		return interfaces.RepositoryCard{}, fmt.Errorf("failed to encrypt cvv: %w", err)
	}

	encryptNotes, err := crypto.Encrypt(card.notes, key)
	if err != nil {
		return interfaces.RepositoryCard{}, fmt.Errorf("failed to encrypt notes: %w", err)
	}

	return interfaces.RepositoryCard{
		ID:                      card.GetID(),
		UserID:                  card.GetUserID(),
		Name:                    card.GetName(),
		EncryptedNumber:         encryptNumber,
		EncryptedExpiredDate:    encryptExpiredDate,
		EncryptedCardHolderName: encryptCardHolderName,
		EncryptedCvv:            encryptCvv,
		EncryptedNotes:          encryptNotes,
		Version:                 card.GetVersion(),
	}, nil
}
