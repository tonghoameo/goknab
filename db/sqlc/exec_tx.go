package db

import (
	"context"
	"fmt"
)

// exec Tx function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q) // exec query
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf(" rx error %v, rb err %v ", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}
