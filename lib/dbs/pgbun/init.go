package pgbun

import (
	"fmt"
	"server/lib/helpers/db/pgbun"
	"sync"
)

var (
	initDBOnce sync.Once
)

func InitDB(appName string, isLogLevelDebug bool) error {
	var errs error
	initDBOnce.Do(func() {
		config, err := DBConfig(appName, isLogLevelDebug)
		if err != nil {
			errs = fmt.Errorf("failed to initalize DB: %v", err)
			return
		}

		DB, err = pgbun.Connect(config)
		if err != nil {
			return
		}
	})

	return errs
}
