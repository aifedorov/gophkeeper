package card

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		id             string
		cardName       string
		number         string
		expiredDate    string
		cardHolderName string
		cvv            string
		notes          string
		version        int64
	}{
		{
			name:           "creates card with all fields",
			id:             testID,
			cardName:       testName,
			number:         testNumber,
			expiredDate:    testExpiredDate,
			cardHolderName: testCardHolderName,
			cvv:            testCvv,
			notes:          testNotes,
			version:        testVersion,
		},
		{
			name:           "creates card without notes",
			id:             "id-2",
			cardName:       "card-2",
			number:         "5500000000000004",
			expiredDate:    "01/26",
			cardHolderName: "Jane Doe",
			cvv:            "456",
			notes:          "",
			version:        2,
		},
		{
			name:           "creates card with zero version",
			id:             "id-3",
			cardName:       "card-3",
			number:         "378282246310005",
			expiredDate:    "06/27",
			cardHolderName: "Bob Smith",
			cvv:            "1234",
			notes:          "amex card",
			version:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			card, err := NewCard(tt.id, tt.cardName, tt.number, tt.expiredDate, tt.cardHolderName, tt.cvv, tt.notes, tt.version)

			require.NoError(t, err)
			require.NotNil(t, card)
			assert.Equal(t, tt.id, card.ID)
			assert.Equal(t, tt.cardName, card.Name)
			assert.Equal(t, tt.number, card.Number)
			assert.Equal(t, tt.expiredDate, card.ExpiredDate)
			assert.Equal(t, tt.cardHolderName, card.CardHolderName)
			assert.Equal(t, tt.cvv, card.Cvv)
			assert.Equal(t, tt.notes, card.Notes)
			assert.Equal(t, tt.version, card.Version)
		})
	}
}

func TestCard_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		card    Card
		wantErr error
	}{
		{
			name: "valid card with all fields",
			card: Card{
				ID:             testID,
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
				Notes:          testNotes,
			},
			wantErr: nil,
		},
		{
			name: "valid card without notes",
			card: Card{
				ID:             testID,
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
				Notes:          "",
			},
			wantErr: nil,
		},
		{
			name: "empty id",
			card: Card{
				ID:             "",
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			wantErr: ErrIDRequired,
		},
		{
			name: "empty name",
			card: Card{
				ID:             testID,
				Name:           "",
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			wantErr: ErrNameRequired,
		},
		{
			name: "empty number",
			card: Card{
				ID:             testID,
				Name:           testName,
				Number:         "",
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			wantErr: ErrNumberRequired,
		},
		{
			name: "empty expired date",
			card: Card{
				ID:             testID,
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    "",
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			wantErr: ErrExpiredDateRequired,
		},
		{
			name: "empty card holder name",
			card: Card{
				ID:             testID,
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: "",
				Cvv:            testCvv,
			},
			wantErr: ErrCardHolderNameRequired,
		},
		{
			name: "empty cvv",
			card: Card{
				ID:             testID,
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            "",
			},
			wantErr: ErrCvvRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.card.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
