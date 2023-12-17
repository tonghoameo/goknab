package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// VerifyEmailTxParams contains the input parameters of the transfer transaction
type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

// VerifyEmailTxResult is the result of the transfer transaction
type VerifyEmailTxResult struct {
	User        User
	verifyEmail VerifyEmail
}

// VerifyEmailTx performs a money transfer from account to other account
// Its creates a transfer record, add account entries and update account's balance within a single databse transaction

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		verifyEmail, err := q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}
		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: verifyEmail.Username,
			//IsEmailVerified: sql.NullBool{
			IsEmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})
		return err
	})

	return result, err
}
