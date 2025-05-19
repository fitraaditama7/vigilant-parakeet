package pgbun

import (
	"encoding/json"
	"fmt"
	"os"
	"server/lib/helpers/db/pgbun"
	"server/lib/helpers/environment"
	"slices"
	"strings"
)

type Config struct {
	Name             string
	CredentialEnvKey string
}

var configs = map[string]Config{
	"DB": {
		Name:             "db",
		CredentialEnvKey: "PG_DB_CREDENTIAL",
	},
}

var DB pgbun.Database

func PostgresDatabaseNames() []string {
	dbNames := []string{}

	for _, config := range configs {
		dbNames = append(dbNames, config.Name)
	}

	return dbNames
}

func PostgresDatabaseConfig(con Config, isLogLevelDebug bool) (pgbun.PGConfig, error) {
	if !IsKnownPostgresDatabase(con.Name) {
		return pgbun.PGConfig{}, fmt.Errorf("unknown postgres database name '%v' - use on of %v", con.Name, PostgresDatabaseNames())
	}

	config := pgbun.PGConfig{
		Username:        "db",
		Password:        "db",
		Host:            "localhost",
		Port:            5432,
		Database:        con.Name,
		IsLogLevelDebug: isLogLevelDebug,
	}

	if configEnv := os.Getenv(con.CredentialEnvKey); configEnv != "" {
		err := json.Unmarshal([]byte(configEnv), &config)
		if err != nil {
			return config, err
		}
	}
	config.Database = environment.GetEnvVar(fmt.Sprintf("PG_%s_DB", strings.ToUpper(con.Name)), config.Database)

	return config, nil
}

func DBConfig(appName string, isLogLevelDebug bool) (pgbun.PGConfig, error) {
	config, err := PostgresDatabaseConfig(configs["DB"], isLogLevelDebug)
	if err != nil {
		return pgbun.PGConfig{}, err
	}

	config.ApplicationName = appName
	return config, nil
}

func IsKnownPostgresDatabase(dbName string) bool {
	return slices.Contains(PostgresDatabaseNames(), dbName)
}
