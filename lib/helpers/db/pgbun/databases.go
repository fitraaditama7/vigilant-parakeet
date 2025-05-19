package pgbun

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

type Database struct {
	config PGConfig
	conn   *bun.DB
}

type PGConfig struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Database        string `json:"dbname"`
	URL             string
	ApplicationName string
	IsLogLevelDebug bool
	WriteTimeout    uint `json:"writeTimeout"`
	ReadTimeout     uint `json:"readTimeout"`
}

func (t PGConfig) GetURL() string {
	if t.URL != "" {
		return t.URL
	}

	addr := t.Host
	if t.Port != 0 {
		addr = fmt.Sprintf("%s:%d", addr, t.Port)
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		t.Username,
		t.Password,
		addr,
		t.Database,
	)
}

func Connect(config PGConfig) (Database, error) {
	if config.URL != "" {
		conn, err := ORMInstance(config.URL, config.IsLogLevelDebug)
		return Database{config: config, conn: conn}, err
	}

	addr := config.Host
	if config.Port != 0 {
		addr = fmt.Sprintf("%s:%d", addr, config.Port)
	}

	conn, err := ORMInstanceFromOptions(&pgdriver.Config{
		Addr:         addr,
		User:         config.Username,
		Password:     config.Password,
		Database:     config.Database,
		AppName:      config.ApplicationName,
		ReadTimeout:  convertToDurationOrDefault(config.ReadTimeout, 30*time.Second),
		WriteTimeout: convertToDurationOrDefault(config.WriteTimeout, 30*time.Second),
	}, config.IsLogLevelDebug)
	if err != nil {
		return Database{}, err
	}

	return Database{config: config, conn: conn}, nil
}

// PGInstance returns the underlying Postgres connection.
// It is not recommended to close this connection directly,
// but the caller to Connect is responsible for closing the connection.
func (t Database) PGInstance() *bun.DB {
	return t.conn
}

func (t Database) Run(callback ORMCallback) error {
	return ORMRun(t, callback)
}

func (t Database) RunContext(ctx context.Context, callback ORMCallbackWithContext) error {
	return ORMRunWithContext(ctx, t, callback)
}

func (t Database) RunInTransaction(callback ORMCallback) error {
	return ORMRunInTransaction(t, callback)
}

func (t Database) RunInTransactionContext(ctx context.Context, callback ORMCallbackWithContext) error {
	return ORMRunInTransactionWithContext(ctx, t, callback)
}

func (t Database) RegisterModel(models ...interface{}) {
	t.conn.RegisterModel(models...)
}

func (t Database) Close() error {
	return t.conn.Close()
}

// Helper generics functions are implemented as standalone functions because
// Go doesn't allow generics in methods
func RunQuery[T any](t Database, callback ORMQueryCallback[T]) (res T, err error) {
	err = t.Run(func(db bun.IDB) error {
		res, err = callback(db)
		return err
	})

	return
}

func RunQueryWithContext[T any](ctx context.Context, t Database, callback ORMQueryCallbackWithContext[T]) (res T, err error) {
	err = t.RunContext(ctx, func(ctx context.Context, db bun.IDB) error {
		res, err = callback(ctx, db)
		return err
	})

	return
}

func RunQueryWithCount[T any](t Database, callback ORMQueryCountCallback[T]) (res T, count int, err error) {
	err = t.Run(func(db bun.IDB) error {
		res, count, err = callback(db)
		return err
	})
	return
}

func RunQueryWithCountContext[T any](ctx context.Context, t Database, callback ORMQueryCountCallbackWithContext[T]) (res T, count int, err error) {
	err = t.RunContext(ctx, func(ctx context.Context, db bun.IDB) error {
		res, count, err = callback(ctx, db)
		return err
	})
	return
}

func RunQueryWithTransaction[T any](t Database, callback ORMQueryCallback[T]) (res T, err error) {
	err = t.RunInTransaction(func(db bun.IDB) error {
		res, err = callback(db)
		return err
	})
	return
}

func RunQueryWithCountInTransaction[T any](t Database, callback ORMQueryCountCallback[T]) (res T, count int, err error) {
	err = t.RunInTransaction(func(db bun.IDB) error {
		res, count, err = callback(db)
		return err
	})
	return
}

// Converts a uint value to time.Duration in seconds, or returns a default value if the input is zero.
func convertToDurationOrDefault(value uint, defaultValue time.Duration) time.Duration {
	// If the value is 0, it indicates that no timeout was explicitly set,
	// so we use the provided default value instead.
	if value == 0 {
		return defaultValue
	}
	return time.Duration(value) * time.Second
}
