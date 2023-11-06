package util

import (
	"fmt"
	"github.com/manticoresoftware/go-sdk/manticore"
	"os"
	"strconv"
)

func GetEnv(env, fallback string) string {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	return v

}

func GetManticoreClient() *manticore.Client {
	cl := manticore.NewClient()
	manticoreHost := GetEnv("MANTICORE_HOST", "127.0.0.1")
	// By default, 9312 listens to both HTTP and SQL requests.
	manticorePort := GetEnv("MANTICORE_PORT", "9312")
	manticorePortNumber, _ := strconv.Atoi(manticorePort)
	cl.SetServer(manticoreHost, uint16(manticorePortNumber))
	_, err := cl.Open()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &cl
}
