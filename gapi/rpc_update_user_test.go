package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	mockdb "github.com/binbomb/goapp/simplebank/db/mock"
	db "github.com/binbomb/goapp/simplebank/db/sqlc"
	"github.com/binbomb/goapp/simplebank/pb"
	"github.com/binbomb/goapp/simplebank/token"
	"github.com/binbomb/goapp/simplebank/utils"
	mockwk "github.com/binbomb/goapp/simplebank/worker/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestUpdateUserAPI(t *testing.T) {
	user, _ := randomUser(t)
	newFullname := utils.RandomString(6)
	newEmail := utils.RandomEmail()
	testCases := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newFullname,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.UpdateUserParams{
					Username: user.Username,
					FullName: sql.NullString{
						String: newFullname,
						Valid:  true,
					},
					Email: sql.NullString{
						String: newEmail,
						Valid:  true,
					},
				}
				updatedUser := db.User{
					Username:          user.Username,
					HashedPassword:    user.HashedPassword,
					FullName:          newFullname,
					Email:             newEmail,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
					IsEmailVerified:   user.IsEmailVerified,
				}
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedUser, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				accessToken, _, err := tokenMaker.CreateToken(user.Username, time.Minute)
				require.NoError(t, err)
				bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
				md := metadata.MD{
					authorizationHeader: []string{
						bearerToken,
					},
				}
				return metadata.NewIncomingContext(context.Background(), md)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				updatedUser := res.GetUser()
				require.Equal(t, user.Username, updatedUser.Username)
				require.Equal(t, newFullname, updatedUser.FullName)
				require.Equal(t, newEmail, updatedUser.Email)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()

			store := mockdb.NewMockStore(storeCtrl)

			//start test server abd send request
			taskCtrl := gomock.NewController(t)
			defer taskCtrl.Finish()
			taskDistributor := mockwk.NewMockTaskDistributor(taskCtrl)

			tc.buildStubs(store)
			server := newTestServer(t, store, taskDistributor)
			ctx := tc.buildContext(t, server.tokenMaker)
			res, err := server.UpdateUser(ctx, tc.req)
			// check response
			tc.checkResponse(t, res, err)
		})
	}
}
