package pgbun

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun/driver/pgdriver"
)

func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, sql.ErrNoRows)
}

type DetailedError struct {
	Err *pgdriver.Error
}

func (e DetailedError) Error() string {
	msg := fmt.Sprintf("ERROR #%v %s", e.Err.Field('C'), e.Err.Field('M'))
	if e.Err.Field('D') != "" {
		msg = fmt.Sprintf("%s, %s", msg, e.Err.Field('D'))
	}

	return msg
}

func (e DetailedError) Unwrap() error {
	return e.Err
}

func (e DetailedError) IsUniqueViolation() bool {
	return EqualsCode(e.Err, UniqueContraintViolation)
}

func AsDetailedError(err error) error {
	var pgErr *pgdriver.Error
	if errors.As(err, &pgErr) {
		return DetailedError{
			Err: pgErr,
		}
	}
	return err
}
