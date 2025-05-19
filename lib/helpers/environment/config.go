package environment

import "os"

func GetEnvVar(key string, defVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defVal
}
