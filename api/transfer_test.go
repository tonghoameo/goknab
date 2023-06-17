package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/binbomb/goapp/simplebank/db/mock"
	db "github.com/binbomb/goapp/simplebank/db/sqlc"
	"github.com/binbomb/goapp/simplebank/utils"
	"github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTransferAPI(t *testing.T) {
	amount := int64(10)
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)
	//	user3, _ := randomUser(t)
	acc1 := randomAccount(user1.Username)
	acc2 := randomAccount(user2.Username)
	//acc3 := randomAccount()

	acc1.Currency = utils.USD
	acc2.Currency = utils.USD
	//acc3.Currency = utils.CAD
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": acc1.ID,
				"to_account_id":   acc2.ID,
				"amount":          amount,
				"currency":        utils.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acc1.ID)).Times(1).Return(acc1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acc2.ID)).Times(1).Return(acc2, nil)
				arg := db.TransferTxParams{
					FromAccountID: acc1.ID,
					ToAccountID:   acc2.ID,
					Amount:        amount,
				}
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/transfers"

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
