package db

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var ErrRecordNotFound = pgx.ErrNoRows

// code error name https://www.postgresql.org/docs/current/errcodes-appendix.html

var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolation,
}

func ErrorDbHandle(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		fmt.Println(pgErr.ConstraintName)
		return pgErr.Code
	}
	return ""
}
