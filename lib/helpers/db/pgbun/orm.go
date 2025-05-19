package pgbun

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

// ORMCallback defines the signature of an external function that will run in
// the context provided by ORMRun.
type ORMCallback func(db bun.IDB) error

// ORMCallbackWithContext defines the signature of an external function that will run in
// the context provided by ORMRun.
// It is used to pass the context and the database connection to the function.
type ORMCallbackWithContext func(ctx context.Context, db bun.IDB) error

// ORMQueryCallback defines the signature of an external function that will run a query in
// the context provided by ORMRun and return a result of a particular type.
type ORMQueryCallback[T any] func(db bun.IDB) (T, error)

// ORMQueryCallbackWithContext defines the signature of an external function that will run a query in
// the context provided by ORMRun and return a result of a particular type.
// It is used to pass the context and the database connection to the function.
type ORMQueryCallbackWithContext[T any] func(ctx context.Context, db bun.IDB) (T, error)

// ORMQueryCountCallback defines the signature of an external function that will run a query in
// the context provided by ORMRun and return a result of a particular type together with count.
type ORMQueryCountCallback[T any] func(db bun.IDB) (T, int, error)

// ORMQueryCountCallbackWithContext defines the signature of an external function that will run a query in
// the context provided by ORMRun and return a result of a particular type together with count.
// It is used to pass the context and the database connection to the function.
type ORMQueryCountCallbackWithContext[T any] func(ctx context.Context, db bun.IDB) (T, int, error)

func ORMInstance(url string, isLogDebug bool) (*bun.DB, error) {
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(url)))
	if pgdb == nil {
		return nil, errors.New("failed to open database connection")
	}

	db := bun.NewDB(pgdb, pgdialect.New())
	if db == nil {
		return nil, errors.New("failed to initiate database connection")
	}

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(isLogDebug)))
	return db, nil
}

func ORMInstanceFromOptions(options *pgdriver.Config, isLogDebug bool) (*bun.DB, error) {
	conn := pgdriver.NewConnector(
		pgdriver.WithAddr(options.Addr),
		pgdriver.WithDatabase(options.Database),
		pgdriver.WithUser(options.User),
		pgdriver.WithPassword(options.Password),
		pgdriver.WithApplicationName(options.AppName),
		pgdriver.WithInsecure(true),
		pgdriver.WithReadTimeout(options.ReadTimeout),
		pgdriver.WithWriteTimeout(options.WriteTimeout),
	)

	pgdb := sql.OpenDB(conn)
	if pgdb == nil {
		return nil, errors.New("failed to open database connection")
	}

	db := bun.NewDB(pgdb, pgdialect.New())
	if db == nil {
		return nil, errors.New("failed to initiate database connection")
	}

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(isLogDebug)))

	return db, db.Ping()
}

// ORMRun provides an ORM execution context to an external function.
func ORMRun(db Database, callback ORMCallback) error {
	return ORMRunWithContext(context.Background(), db, func(_ context.Context, db bun.IDB) error {
		return callback(db)
	})
}

// ORMRunWithContext provides an ORM execution context to an external function.
// It is used to pass the context and the database connection to the function.
// The context is passed to the function, and the database connection is
// passed as a bun.IDB interface.
func ORMRunWithContext(ctx context.Context, db Database, callback ORMCallbackWithContext) error {
	err := callback(ctx, db.conn)
	if err != nil && db.config.IsLogLevelDebug {
		fmt.Println("Postgres Transaction error:")
		jsObj, _ := json.MarshalIndent(err, "", "\t")
		fmt.Printf("%v\n", fmt.Sprintln(string(jsObj)))
	}

	return err
}

// ORMRunInTransaction provides an transactioned ORM execution context to an
// external function.
func ORMRunInTransaction(db Database, callback ORMCallback) error {
	return ORMRunInTransactionWithContext(context.Background(), db, func(_ context.Context, tx bun.IDB) error {
		return callback(tx)
	})
}

// ORMRunInTransactionWithContext provides an transactioned ORM execution context to an
// external function.
// It is used to pass the context and the database connection to the function.
// The context is passed to the function, and the database connection is
// passed as a bun.IDB interface.
func ORMRunInTransactionWithContext(ctx context.Context, db Database, callback ORMCallbackWithContext) error {
	err := db.conn.RunInTx(ctx, &sql.TxOptions{}, func(_ context.Context, tx bun.Tx) error {
		return callback(ctx, tx)
	})

	if err != nil && db.config.IsLogLevelDebug {
		fmt.Println("Postgres Transaction error:")
		jsObj, _ := json.MarshalIndent(err, "", "\t")
		fmt.Printf("%v\n", fmt.Sprintln(string(jsObj)))
	}

	return err
}
