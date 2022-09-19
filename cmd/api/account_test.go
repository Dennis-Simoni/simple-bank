package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
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

	server := NewServer(store)
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(a.ID)).Times(1).Return(a, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", a.ID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}
