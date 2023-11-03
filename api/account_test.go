package api

import (
	mockdb "simple_bank/db/mock"
	db "simple_bank/db/sqlc"
	"simple_bank/util"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {
	ac := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	// Build stubs
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(ac.ID)).
		Times(1).
		Return(ac, nil)
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
