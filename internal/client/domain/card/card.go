// Package card provides card domain entities for the GophKeeper client.
package card

// Card represents a card entity in the client domain.
// It contains card information (number, expiration date, CVV, etc.) along with metadata.
type Card struct {
	ID             string // Unique card identifier
	Name           string // Display name (e.g., "Visa Main Card")
	Number         string // Card number
	ExpiredDate    string // Expiration date
	CardHolderName string // Cardholder name
	Cvv            string // CVV code
	Notes          string // Optional notes/metadata
	Version        int64  // Version number
}

// NewCard creates a new Card entity with the provided data.
// Returns an error if validation fails (currently always returns nil as validation is done in Validate method).
func NewCard(id, name, number, expiredDate, cardHolderName, cvv, notes string, version int64) (*Card, error) {
	return &Card{
		ID:             id,
		Name:           name,
		Number:         number,
		ExpiredDate:    expiredDate,
		CardHolderName: cardHolderName,
		Cvv:            cvv,
		Notes:          notes,
		Version:        version,
	}, nil
}

// Validate checks that all required fields (ID, Name, Number, ExpiredDate, CardHolderName, Cvv) are not empty.
// Returns an error if any required field is missing.
func (c *Card) Validate() error {
	if c.ID == "" {
		return ErrIDRequired
	}
	if c.Name == "" {
		return ErrNameRequired
	}
	if c.Number == "" {
		return ErrNumberRequired
	}
	if c.ExpiredDate == "" {
		return ErrExpiredDateRequired
	}
	if c.CardHolderName == "" {
		return ErrCardHolderNameRequired
	}
	if c.Cvv == "" {
		return ErrCvvRequired
	}
	return nil
}
