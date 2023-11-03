package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	mockdb "simple_bank/db/mock"
	db "simple_bank/db/sqlc"
	"simple_bank/util"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	ac := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	// Build stubs.
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(ac.ID)).
		// Times(1) means we expect this function to be called exactly 1 time.
		Times(1).
		// We can use the Return function to tell gomock to return some specific values
		// whenever the GetAccount function is called.
		Return(ac, nil)

	// Start test server and send request.
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", ac.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	// Check the response.
	require.Equal(t, http.StatusOK, recorder.Code)
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
