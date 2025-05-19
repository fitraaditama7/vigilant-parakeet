package pgbun

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/schema"
)

func inArray(b string, vals []string) bool {
	for _, a := range vals {
		if a == b {
			return true
		}
	}
	return false
}

func EqualsCode(err error, codes ...string) bool {
	var pgErr *pgdriver.Error
	if !errors.As(err, &pgErr) {
		return false
	}
	return inArray(pgErr.Field('C'), codes)
}

func StringsToGeneric(items []string) []interface{} {
	res := make([]interface{}, 0, len(items))
	for _, item := range items {
		res = append(res, item)
	}

	return res
}

func UUIDsToGeneric(items []uuid.UUID) []interface{} {
	res := make([]interface{}, 0, len(items))
	for _, item := range items {
		res = append(res, item)
	}

	return res
}

func InMulti(items ...interface{}) schema.QueryAppender {
	return bun.In(items)
}
