package config

import "os"

func IsProduction() bool {
	return os.Getenv("ENV") == "production"
}
