package card

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	s := newTestSetup(t)
	defer s.cleanup()
	service := NewService(s.mockClient, s.mockCache)

	require.NotNil(t, service)
}

func TestService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		card    Card
		setup   func(*testSetup, Card)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful creation",
			card: Card{
				ID:             testID,
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
				Notes:          testNotes,
			},
			setup: func(s *testSetup, card Card) {
				s.expectCreateSuccess(card.ID, testVersion)
			},
			wantErr: false,
		},
		{
			name: "successful creation without notes",
			card: Card{
				ID:             "test-id-456",
				Name:           "another-card",
				Number:         "5500000000000004",
				ExpiredDate:    "01/26",
				CardHolderName: "Jane Doe",
				Cvv:            "456",
				Notes:          "",
			},
			setup: func(s *testSetup, card Card) {
				s.expectCreateSuccess(card.ID, testVersion)
			},
			wantErr: false,
		},
		{
			name: "validation error - empty id",
			card: Card{
				ID:             "",
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			setup:   func(s *testSetup, card Card) {},
			wantErr: true,
			errMsg:  "invalid card",
		},
		{
			name: "validation error - empty name",
			card: Card{
				ID:             testID,
				Name:           "",
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			setup:   func(s *testSetup, card Card) {},
			wantErr: true,
			errMsg:  "invalid card",
		},
		{
			name: "validation error - empty number",
			card: Card{
				ID:             testID,
				Name:           testName,
				Number:         "",
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			setup:   func(s *testSetup, card Card) {},
			wantErr: true,
			errMsg:  "invalid card",
		},
		{
			name: "validation error - empty expired date",
			card: Card{
				ID:             testID,
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    "",
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			setup:   func(s *testSetup, card Card) {},
			wantErr: true,
			errMsg:  "invalid card",
		},
		{
			name: "validation error - empty card holder name",
			card: Card{
				ID:             testID,
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: "",
				Cvv:            testCvv,
			},
			setup:   func(s *testSetup, card Card) {},
			wantErr: true,
			errMsg:  "invalid card",
		},
		{
			name: "validation error - empty cvv",
			card: Card{
				ID:             testID,
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            "",
			},
			setup:   func(s *testSetup, card Card) {},
			wantErr: true,
			errMsg:  "invalid card",
		},
		{
			name: "client error",
			card: Card{
				ID:             testID,
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			setup: func(s *testSetup, card Card) {
				s.expectCreateClientError(errors.New("network error"))
			},
			wantErr: true,
			errMsg:  "failed to create card",
		},
		{
			name: "cache error on set version",
			card: Card{
				ID:             "cache-error-id",
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			setup: func(s *testSetup, card Card) {
				s.mockClient.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return("cache-error-id", int64(1), nil).
					Times(1)
				s.mockCache.EXPECT().
					SetCardVersion("cache-error-id", int64(1)).
					Return(errors.New("cache error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to save card to cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s, tt.card)

			err := s.service.Create(s.ctx, tt.card)

			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}

func TestService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(*testSetup)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful list with multiple cards",
			setup: func(s *testSetup) {
				cards := []Card{
					{ID: "id-1", Name: "card-1", Number: "4111111111111111", ExpiredDate: "12/25", CardHolderName: "John", Cvv: "123", Version: 1},
					{ID: "id-2", Name: "card-2", Number: "5500000000000004", ExpiredDate: "01/26", CardHolderName: "Jane", Cvv: "456", Version: 1},
				}
				s.expectListSuccess(cards)
				s.wantCards = cards
			},
			wantErr: false,
		},
		{
			name: "successful list with empty result",
			setup: func(s *testSetup) {
				s.expectListSuccess([]Card{})
				s.wantCards = []Card{}
			},
			wantErr: false,
		},
		{
			name: "successful list with single card",
			setup: func(s *testSetup) {
				cards := []Card{
					{ID: "single-id", Name: "single-card", Number: "4111111111111111", ExpiredDate: "12/25", CardHolderName: "John", Cvv: "123", Version: 2},
				}
				s.expectListSuccess(cards)
				s.wantCards = cards
			},
			wantErr: false,
		},
		{
			name: "client error",
			setup: func(s *testSetup) {
				s.expectListError(errors.New("network error"))
			},
			wantErr: true,
			errMsg:  "failed to get list of cards",
		},
		{
			name: "server returns invalid version 0",
			setup: func(s *testSetup) {
				cards := []Card{
					{ID: "id-1", Name: "card-1", Number: "4111111111111111", ExpiredDate: "12/25", CardHolderName: "John", Cvv: "123", Version: 0},
				}
				s.mockClient.EXPECT().
					List(gomock.Any()).
					Return(cards, nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "server returned invalid version 0",
		},
		{
			name: "cache error on set version",
			setup: func(s *testSetup) {
				cards := []Card{
					{ID: "id-1", Name: "card-1", Number: "4111111111111111", ExpiredDate: "12/25", CardHolderName: "John", Cvv: "123", Version: 1},
				}
				s.mockClient.EXPECT().
					List(gomock.Any()).
					Return(cards, nil).
					Times(1)
				s.mockCache.EXPECT().
					SetCardVersion("id-1", int64(1)).
					Return(errors.New("cache error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to save card to cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s)

			cards, err := s.service.List(s.ctx)

			assertError(t, err, tt.wantErr, tt.errMsg)
			if !tt.wantErr && s.wantCards != nil {
				assertCardsEqual(t, cards, s.wantCards)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		card    Card
		setup   func(*testSetup, string, Card)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful update",
			id:   testID,
			card: Card{
				ID:             testID,
				Name:           "updated-card",
				Number:         "4111111111111111",
				ExpiredDate:    "12/26",
				CardHolderName: "Updated Name",
				Cvv:            "789",
				Notes:          "updated notes",
			},
			setup: func(s *testSetup, id string, card Card) {
				s.expectUpdateSuccess(id, testVersion, int64(2))
			},
			wantErr: false,
		},
		{
			name: "successful update without notes",
			id:   "test-id-456",
			card: Card{
				ID:             "test-id-456",
				Name:           "updated-card",
				Number:         "5500000000000004",
				ExpiredDate:    "01/27",
				CardHolderName: "Jane Updated",
				Cvv:            "999",
				Notes:          "",
			},
			setup: func(s *testSetup, id string, card Card) {
				s.expectUpdateSuccess(id, int64(3), int64(4))
			},
			wantErr: false,
		},
		{
			name: "validation error - empty name",
			id:   "test-id-789",
			card: Card{
				ID:             "test-id-789",
				Name:           "",
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			setup:   func(s *testSetup, id string, card Card) {},
			wantErr: true,
			errMsg:  "invalid card",
		},
		{
			name: "validation error - empty number",
			id:   "test-id-101",
			card: Card{
				ID:             "test-id-101",
				Name:           testName,
				Number:         "",
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			setup:   func(s *testSetup, id string, card Card) {},
			wantErr: true,
			errMsg:  "invalid card",
		},
		{
			name: "client error - not found",
			id:   "non-existent-id",
			card: Card{
				ID:             "non-existent-id",
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			setup: func(s *testSetup, id string, card Card) {
				s.expectUpdateClientError(id, testVersion, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "failed to update",
		},
		{
			name: "client error - network error",
			id:   "test-id-103",
			card: Card{
				ID:             "test-id-103",
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			setup: func(s *testSetup, id string, card Card) {
				s.expectUpdateClientError(id, int64(2), errors.New("network error"))
			},
			wantErr: true,
			errMsg:  "failed to update",
		},
		{
			name: "cache error - get version fails",
			id:   "test-id-104",
			card: Card{
				ID:             "test-id-104",
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			setup: func(s *testSetup, id string, card Card) {
				s.mockCache.EXPECT().
					GetCardVersion(id).
					Return(int64(0), errors.New("cache miss")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to get version from cache",
		},
		{
			name: "cache error - set version fails after update",
			id:   "test-id-105",
			card: Card{
				ID:             "test-id-105",
				Name:           testName,
				Number:         testNumber,
				ExpiredDate:    testExpiredDate,
				CardHolderName: testCardHolderName,
				Cvv:            testCvv,
			},
			setup: func(s *testSetup, id string, card Card) {
				s.mockCache.EXPECT().
					GetCardVersion(id).
					Return(int64(1), nil).
					Times(1)
				s.mockClient.EXPECT().
					Update(gomock.Any(), id, gomock.Any()).
					Return(int64(2), nil).
					Times(1)
				s.mockCache.EXPECT().
					SetCardVersion(gomock.Any(), int64(2)).
					Return(errors.New("cache error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to save card to cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s, tt.id, tt.card)

			err := s.service.Update(s.ctx, tt.id, tt.card)

			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}

func TestService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		setup   func(*testSetup, string)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful deletion",
			id:   testID,
			setup: func(s *testSetup, id string) {
				s.expectDeleteSuccess(id)
			},
			wantErr: false,
		},
		{
			name: "successful deletion with different ID",
			id:   "another-id-456",
			setup: func(s *testSetup, id string) {
				s.expectDeleteSuccess(id)
			},
			wantErr: false,
		},
		{
			name: "client error - not found",
			id:   "non-existent-id",
			setup: func(s *testSetup, id string) {
				s.expectDeleteClientError(id, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "failed to delete card",
		},
		{
			name: "client error - network error",
			id:   "test-id-789",
			setup: func(s *testSetup, id string) {
				s.expectDeleteClientError(id, errors.New("network error"))
			},
			wantErr: true,
			errMsg:  "failed to delete card",
		},
		{
			name: "cache error - delete version fails",
			id:   "test-id-999",
			setup: func(s *testSetup, id string) {
				s.mockClient.EXPECT().
					Delete(gomock.Any(), id).
					Return(nil).
					Times(1)
				s.mockCache.EXPECT().
					DeleteCardVersion(id).
					Return(errors.New("cache error")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to delete card from cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := newTestSetup(t)
			defer s.cleanup()
			s.initService()
			tt.setup(s, tt.id)

			err := s.service.Delete(s.ctx, tt.id)

			assertError(t, err, tt.wantErr, tt.errMsg)
		})
	}
}
