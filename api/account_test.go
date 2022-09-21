package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"testing"
	"time"
)

func TestCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := db.Account{
		ID:        1,
		Owner:     "bob",
		Balance:   0,
		Currency:  "USD",
		CreatedAt: time.Now(),
	}

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().CreateAccount(gomock.Any(), db.CreateAccountParams{
		Owner:    a.Owner,
		Currency: a.Currency,
	}).Return(a, nil)

	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	reqs := createAccountRequest{
		Owner:    a.Owner,
		Currency: a.Currency,
	}

	b, err := json.Marshal(reqs)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(b))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestGetAccount(t *testing.T) {

	a := randomAccount()

	cases := []struct {
		name          string
		accountID     int64
		buildStubs    func(s *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: a.ID,
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().GetAccount(gomock.Any(), gomock.Eq(a.ID)).Times(1).Return(a, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, a)
			},
		},
		{
			name:      "Not found",
			accountID: a.ID,
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().GetAccount(gomock.Any(), gomock.Eq(a.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal error",
			accountID: a.ID,
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().GetAccount(gomock.Any(), gomock.Eq(a.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Invalid account id",
			accountID: 0,
			buildStubs: func(s *mockdb.MockStore) {
				s.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tt.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tt.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
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
