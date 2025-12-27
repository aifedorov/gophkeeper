package card

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testID             = "test-id-123"
	testName           = "test-card"
	testNumber         = "4111111111111111"
	testExpiredDate    = "12/25"
	testCardHolderName = "John Doe"
	testCvv            = "123"
	testNotes          = "test notes"
	testVersion        = int64(1)
)

type testSetup struct {
	ctrl       *gomock.Controller
	mockClient *MockClient
	mockCache  *MockCacheStorage
	service    Service
	ctx        context.Context
	testCard   Card
	wantCards  []Card
}

func newTestSetup(t *testing.T) *testSetup {
	ctrl := gomock.NewController(t)

	return &testSetup{
		ctrl:       ctrl,
		mockClient: NewMockClient(ctrl),
		mockCache:  NewMockCacheStorage(ctrl),
		ctx:        context.Background(),
		testCard: Card{
			ID:             testID,
			Name:           testName,
			Number:         testNumber,
			ExpiredDate:    testExpiredDate,
			CardHolderName: testCardHolderName,
			Cvv:            testCvv,
			Notes:          testNotes,
		},
	}
}

func (s *testSetup) initService() {
	s.service = NewService(s.mockClient, s.mockCache)
}

func (s *testSetup) cleanup() {
	s.ctrl.Finish()
}

func (s *testSetup) expectCreateSuccess(id string, version int64) {
	s.mockClient.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(id, version, nil).
		Times(1)
	s.mockCache.EXPECT().
		SetCardVersion(id, version).
		Return(nil).
		Times(1)
}

func (s *testSetup) expectCreateClientError(err error) {
	s.mockClient.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return("", int64(0), err).
		Times(1)
}

func (s *testSetup) expectListSuccess(cards []Card) {
	s.mockClient.EXPECT().
		List(gomock.Any()).
		Return(cards, nil).
		Times(1)
	for _, card := range cards {
		if card.Version > 0 {
			s.mockCache.EXPECT().
				SetCardVersion(card.ID, card.Version).
				Return(nil).
				Times(1)
		}
	}
}

func (s *testSetup) expectListError(err error) {
	s.mockClient.EXPECT().
		List(gomock.Any()).
		Return(nil, err).
		Times(1)
}

func (s *testSetup) expectUpdateSuccess(id string, currentVersion, newVersion int64) {
	s.mockCache.EXPECT().
		GetCardVersion(id).
		Return(currentVersion, nil).
		Times(1)
	s.mockClient.EXPECT().
		Update(gomock.Any(), id, gomock.Any()).
		Return(newVersion, nil).
		Times(1)
	s.mockCache.EXPECT().
		SetCardVersion(gomock.Any(), newVersion).
		Return(nil).
		Times(1)
}

func (s *testSetup) expectUpdateClientError(id string, currentVersion int64, err error) {
	s.mockCache.EXPECT().
		GetCardVersion(id).
		Return(currentVersion, nil).
		Times(1)
	s.mockClient.EXPECT().
		Update(gomock.Any(), id, gomock.Any()).
		Return(int64(0), err).
		Times(1)
}

func (s *testSetup) expectDeleteSuccess(id string) {
	s.mockClient.EXPECT().
		Delete(gomock.Any(), id).
		Return(nil).
		Times(1)
	s.mockCache.EXPECT().
		DeleteCardVersion(id).
		Return(nil).
		Times(1)
}

func (s *testSetup) expectDeleteClientError(id string, err error) {
	s.mockClient.EXPECT().
		Delete(gomock.Any(), id).
		Return(err).
		Times(1)
}

func assertError(t *testing.T, err error, wantErr bool, errMsg string) {
	t.Helper()
	if wantErr {
		require.Error(t, err)
		if errMsg != "" {
			assert.Contains(t, err.Error(), errMsg)
		}
	} else {
		require.NoError(t, err)
	}
}

func assertCardsEqual(t *testing.T, got, want []Card) {
	t.Helper()
	assert.Equal(t, want, got)
}
