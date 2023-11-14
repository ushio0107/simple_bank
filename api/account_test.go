package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
	user, _ := randomUser(t)
	ac := randomAccount(user.Username)

	// Define all test cases to implement 100% test coverage.
	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		/*
			{
				name:      "Example",
				accountID: ac.ID,
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetAccount(gomock.Any(), gomock.Eq(ac.ID)).
						// Times(1) means we expect this function to be called exactly 1 time.
						Times(1).
						// We can use the Return function to tell gomock to return some specific values
						// whenever the GetAccount function is called.
						Return(ac, nil)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					// Check the response.
					require.Equal(t, http.StatusOK, recorder.Code)
					// Check the response body.
					requireBodyMatchAccount(t, recorder.Body, ac)
				},
			}
		*/
		{
			name:      "OK", // First happy case.
			accountID: ac.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(ac.ID)).
					Times(1).
					// Return correct account.
					Return(ac, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, ac)
			},
		}, {
			name:      "NotFound",
			accountID: ac.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(ac.ID)).
					Times(1).
					// ac is empty since the specific account id isn't be found.
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		}, {
			name:      "InternalError",
			accountID: ac.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(ac.ID)).
					Times(1).
					// DB connction is failed/lost.
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		}, {
			name: "InvalidID",
			// The minimum account id is 1.
			accountID: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					// Since the ID is invalid, this GetAccount function should not be called by the handler.
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// Start test server and send request.
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

// func TestListAccountsAPI(t *testing.T) {
// 	user, _ := randomUser(t)

// 	n := 5
// 	acs := make([]db.Account, n)
// 	for i := 0; i < n; i++ {
// 		acs[i] = randomAccount(user.Username)
// 	}

// 	type Query struct {
// 		pageID   int
// 		pageSize int
// 	}

// 	// Define all test cases to implement 100% test coverage.
// 	testCases := []struct {
// 		name          string
// 		query         Query
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK", // First happy case.
// 			query: Query{
// 				pageID:   1,
// 				pageSize: n,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				arg := db.ListAccountsParams{
// 					Owner:  user.Username,
// 					Limit:  int32(n),
// 					Offset: 0,
// 				}
// 				store.EXPECT().
// 					ListAccounts(gomock.Any(), gomock.Eq(arg)).
// 					Times(1).
// 					// Return correct account.
// 					Return(acs, nil)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				requireBodyMatchAccounts(t, recorder.Body, acs)
// 			},
// 		}, {
// 			name: "InternalError",
// 			query: Query{
// 				pageID:   1,
// 				pageSize: n,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					ListAccounts(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					// DB connction is failed/lost.
// 					Return([]db.Account{}, sql.ErrConnDone)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		}, {
// 			name: "InvalidPageID",
// 			// The minimum page id is 1.
// 			query: Query{
// 				pageID:   -1,
// 				pageSize: n,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					ListAccounts(gomock.Any(), gomock.Any()).
// 					// Since the ID is invalid, this ListAccounts function should not be called by the handler.
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		}, {
// 			name: "InvalidPageSize",
// 			query: Query{
// 				pageID:   1,
// 				pageSize: 100000,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					ListAccounts(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			// Start test server and send request.
// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			url := "/accounts"
// 			request, err := http.NewRequest(http.MethodGet, url, nil)
// 			require.NoError(t, err)

// 			// Add query parameters to request URL
// 			q := request.URL.Query()
// 			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
// 			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
// 			request.URL.RawQuery = q.Encode()

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(recorder)
// 		})
// 	}
// }

func randomAccount(username string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
