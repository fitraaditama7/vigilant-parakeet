package pgbun

import "server/lib/helpers/db/pgbun"

type RegisterModelFunc func(database pgbun.Database)

func RegisterModel(fn RegisterModelFunc) {
	if DB.PGInstance() != nil {
		fn(DB)
	}
}
