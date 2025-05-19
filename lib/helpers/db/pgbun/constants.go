package pgbun

const UniqueContraintViolation = "23505"

type QueryOp string

const (
	Eq    QueryOp = "="
	NotEq QueryOp = "!="
	Lt    QueryOp = "<"
	Lte   QueryOp = "<="
	Gt    QueryOp = ">"
	Gte   QueryOp = ">="
)
