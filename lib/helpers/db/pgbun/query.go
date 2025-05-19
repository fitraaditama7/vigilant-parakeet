package pgbun

import (
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Query interface{}

type QueryHook func(q bun.QueryBuilder) bun.QueryBuilder

var EmptyHook = func(q bun.QueryBuilder) bun.QueryBuilder {
	return q
}

func SelectQueryHook(hook func(query *bun.SelectQuery) *bun.SelectQuery) QueryHook {
	return func(q bun.QueryBuilder) bun.QueryBuilder {
		if selectQuery, ok := q.Unwrap().(*bun.SelectQuery); ok {
			return hook(selectQuery).QueryBuilder()
		}
		return q
	}
}

func ApplyQueryBuilderHooks(q bun.QueryBuilder, hooks ...QueryHook) bun.QueryBuilder {
	for _, hookFn := range hooks {
		q = hookFn(q)
	}
	return q
}

func ApplyQueryHooks(hooks ...QueryHook) QueryHook {
	return func(q bun.QueryBuilder) bun.QueryBuilder {
		for _, hookFn := range hooks {
			q = hookFn(q)
		}
		return q
	}
}

func WhereHook[T any](field bun.Ident, op QueryOp, value T) QueryHook {
	return func(q bun.QueryBuilder) bun.QueryBuilder {
		return q.Where(fmt.Sprintf("? %s ?", op), field, value)
	}
}

func WhereInHook(field string, options []interface{}) QueryHook {
	return func(q bun.QueryBuilder) bun.QueryBuilder {
		if len(options) > 0 {
			return q.Where(fmt.Sprintf("%s IN (?)", field), bun.In(options))
		}
		return q.Where("FALSE")
	}
}

func WhereInHookT[T any](field bun.Ident, options []T) QueryHook {
	return func(q bun.QueryBuilder) bun.QueryBuilder {
		if len(options) > 0 {
			return q.Where(fmt.Sprintf("%s IN (?)", field), bun.In(options))
		}
		return q.Where("FALSE")
	}
}

func WhereNotInHook[T any](field bun.Ident, options []T) QueryHook {
	return func(q bun.QueryBuilder) bun.QueryBuilder {
		if len(options) > 0 {
			return q.Where(fmt.Sprintf("? NOT IN (?)", field, bun.In(options)))
		}
		return q.Where("FALSE")
	}
}

func WhereOverlapHookT[T any](field bun.Ident, options []T) QueryHook {
	return func(q bun.QueryBuilder) bun.QueryBuilder {
		if len(options) > 0 {
			return q.Where(fmt.Sprintf("? && (?)", field, bun.In(options)))
		}
		return q.Where("FALSE")
	}
}

func WhereOrInHook(field string, options []interface{}) QueryHook {
	return func(q bun.QueryBuilder) bun.QueryBuilder {
		if len(options) > 0 {
			return q.WhereOr(fmt.Sprintf("%s IN (?)", field), bun.In(options))
		}
		return q.Where("FALSE")
	}
}

func WhereOrOverlapHook[T any](field string, options []interface{}) QueryHook {
	return func(q bun.QueryBuilder) bun.QueryBuilder {
		if len(options) > 0 {
			return q.WhereOr(fmt.Sprintf("%v && ?", field), pgdialect.Array(options))
		}
		return q.Where("FALSE")
	}
}

func WhereOrOverlapHookT[T any](field string, options []T) QueryHook {
	return func(q bun.QueryBuilder) bun.QueryBuilder {
		if len(options) > 0 {
			return q.WhereOr(fmt.Sprintf("%v && ?", field), pgdialect.Array(options))
		}
		return q.Where("FALSE")
	}
}

func PaginateHook(offset, limit int) QueryHook {
	return SelectQueryHook(func(q *bun.SelectQuery) *bun.SelectQuery {
		if limit == 0 {
			return q
		}

		return q.Offset(offset).Limit(limit)
	})
}

func SortHook(sortKey, sortOrder string) QueryHook {
	return SelectQueryHook(func(q *bun.SelectQuery) *bun.SelectQuery {
		if sortKey != "" {
			q = q.OrderExpr(fmt.Sprintf("%s %s", sortKey, sortOrder))
		}
		return q
	})
}

func Paginate(page int, perPage int) QueryHook {
	offset := (page - 1) * perPage
	limit := perPage
	return PaginateHook(offset, limit)
}
