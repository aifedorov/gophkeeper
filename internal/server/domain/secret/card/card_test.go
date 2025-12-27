package card

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testIDCard             = "test-id-123-card"
	testNameCard           = "test-card-name"
	testNumberCard         = "4111111111111111"
	testExpiredDateCard    = "12/25"
	testCardHolderNameCard = "John Doe"
	testCvvCard            = "123"
	testNotesCard          = "test notes card"
	testVersionCard        = int64(1)
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
		wantErr        error
		wantErrStr     string
	}{
		{
			name:           "creates card with all fields",
			id:             testIDCard,
			cardName:       testNameCard,
			number:         testNumberCard,
			expiredDate:    testExpiredDateCard,
			cardHolderName: testCardHolderNameCard,
			cvv:            testCvvCard,
			notes:          testNotesCard,
			version:        testVersionCard,
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
			name:           "empty id",
			id:             "",
			cardName:       testNameCard,
			number:         testNumberCard,
			expiredDate:    testExpiredDateCard,
			cardHolderName: testCardHolderNameCard,
			cvv:            testCvvCard,
			notes:          testNotesCard,
			version:        testVersionCard,
			wantErr:        ErrIDRequired,
		},
		{
			name:           "empty name",
			id:             testIDCard,
			cardName:       "",
			number:         testNumberCard,
			expiredDate:    testExpiredDateCard,
			cardHolderName: testCardHolderNameCard,
			cvv:            testCvvCard,
			notes:          testNotesCard,
			version:        testVersionCard,
			wantErr:        ErrNameRequired,
		},
		{
			name:           "empty number",
			id:             testIDCard,
			cardName:       testNameCard,
			number:         "",
			expiredDate:    testExpiredDateCard,
			cardHolderName: testCardHolderNameCard,
			cvv:            testCvvCard,
			notes:          testNotesCard,
			version:        testVersionCard,
			wantErr:        ErrNumberRequired,
		},
		{
			name:           "empty expired date",
			id:             testIDCard,
			cardName:       testNameCard,
			number:         testNumberCard,
			expiredDate:    "",
			cardHolderName: testCardHolderNameCard,
			cvv:            testCvvCard,
			notes:          testNotesCard,
			version:        testVersionCard,
			wantErr:        ErrExpiredDateRequired,
		},
		{
			name:           "empty card holder name",
			id:             testIDCard,
			cardName:       testNameCard,
			number:         testNumberCard,
			expiredDate:    testExpiredDateCard,
			cardHolderName: "",
			cvv:            testCvvCard,
			notes:          testNotesCard,
			version:        testVersionCard,
			wantErr:        ErrCardHolderNameRequired,
		},
		{
			name:           "empty cvv",
			id:             testIDCard,
			cardName:       testNameCard,
			number:         testNumberCard,
			expiredDate:    testExpiredDateCard,
			cardHolderName: testCardHolderNameCard,
			cvv:            "",
			notes:          testNotesCard,
			version:        testVersionCard,
			wantErr:        ErrCvvRequired,
		},
		{
			name:           "zero version",
			id:             testIDCard,
			cardName:       testNameCard,
			number:         testNumberCard,
			expiredDate:    testExpiredDateCard,
			cardHolderName: testCardHolderNameCard,
			cvv:            testCvvCard,
			notes:          testNotesCard,
			version:        0,
			wantErrStr:     "invalid card version",
		},
		{
			name:           "negative version",
			id:             testIDCard,
			cardName:       testNameCard,
			number:         testNumberCard,
			expiredDate:    testExpiredDateCard,
			cardHolderName: testCardHolderNameCard,
			cvv:            testCvvCard,
			notes:          testNotesCard,
			version:        -1,
			wantErrStr:     "invalid card version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			card, err := NewCard(tt.id, tt.cardName, tt.number, tt.expiredDate, tt.cardHolderName, tt.cvv, tt.notes, tt.version)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, card)
				return
			}
			if tt.wantErrStr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrStr)
				assert.Nil(t, card)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, card)
			assert.Equal(t, tt.id, card.GetID())
			assert.Equal(t, tt.cardName, card.GetName())
			assert.Equal(t, tt.number, card.GetNumber())
			assert.Equal(t, tt.expiredDate, card.GetExpiredDate())
			assert.Equal(t, tt.cardHolderName, card.GetCardHolderName())
			assert.Equal(t, tt.cvv, card.GetCvv())
			assert.Equal(t, tt.notes, card.GetNotes())
			assert.Equal(t, tt.version, card.GetVersion())
		})
	}
}

func TestCard_Getters(t *testing.T) {
	t.Parallel()

	card, err := NewCard(testIDCard, testNameCard, testNumberCard, testExpiredDateCard, testCardHolderNameCard, testCvvCard, testNotesCard, testVersionCard)
	require.NoError(t, err)

	assert.Equal(t, testIDCard, card.GetID())
	assert.Equal(t, "", card.GetUserID()) // userID is not set by NewCard
	assert.Equal(t, testNameCard, card.GetName())
	assert.Equal(t, testNumberCard, card.GetNumber())
	assert.Equal(t, testExpiredDateCard, card.GetExpiredDate())
	assert.Equal(t, testCardHolderNameCard, card.GetCardHolderName())
	assert.Equal(t, testCvvCard, card.GetCvv())
	assert.Equal(t, testNotesCard, card.GetNotes())
	assert.Equal(t, testVersionCard, card.GetVersion())
}
