package card

import "fmt"

// Card represents a payment card entity in the card domain.
// It contains card information (number, expiration date, CVV, etc.) along with metadata,
// all of which are encrypted before storage.
type Card struct {
	id             string
	userID         string
	name           string
	number         string
	expiredDate    string
	cardHolderName string
	cvv            string
	notes          string
	version        int64
}

// NewCard creates a new Card entity with the provided data.
// It validates that all required fields (id, name, number, expiredDate, cardHolderName, cvv) are not empty.
// Returns an error if validation fails.
func NewCard(id, name, number, expiredDate, cardHolderName, cvv, notes string, version int64) (*Card, error) {
	if id == "" {
		return nil, ErrIDRequired
	}
	if name == "" {
		return nil, ErrNameRequired
	}
	if number == "" {
		return nil, ErrNumberRequired
	}
	if expiredDate == "" {
		return nil, ErrExpiredDateRequired
	}
	if cardHolderName == "" {
		return nil, ErrCardHolderNameRequired
	}
	if cvv == "" {
		return nil, ErrCvvRequired
	}
	if version < 1 {
		return nil, fmt.Errorf("invalid card version: %d", version)
	}

	return &Card{
		id:             id,
		name:           name,
		number:         number,
		expiredDate:    expiredDate,
		cardHolderName: cardHolderName,
		cvv:            cvv,
		notes:          notes,
		version:        version,
	}, nil
}

// GetID returns the card's unique identifier.
func (c *Card) GetID() string {
	return c.id
}

// GetUserID returns the ID of the user who owns this card.
func (c *Card) GetUserID() string {
	return c.userID
}

// GetName returns the card's display name for this card.
func (c *Card) GetName() string {
	return c.name
}

// GetNumber returns the decrypted card number for this card.
func (c *Card) GetNumber() string {
	return c.number
}

// GetExpiredDate returns the decrypted expiration date for this card.
func (c *Card) GetExpiredDate() string {
	return c.expiredDate
}

// GetCardHolderName returns the decrypted cardholder name for this card.
func (c *Card) GetCardHolderName() string {
	return c.cardHolderName
}

// GetCvv returns the decrypted CVV for this card.
func (c *Card) GetCvv() string {
	return c.cvv
}

// GetNotes returns the decrypted metadata/notes for this card.
func (c *Card) GetNotes() string {
	return c.notes
}

// GetVersion returns the card's version number for optimistic locking.
func (c *Card) GetVersion() int64 {
	return c.version
}
