package util

import (
	"os"
)

func GetEnv(env, fallback string) string {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	return v

}
