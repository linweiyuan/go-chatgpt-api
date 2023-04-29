package env

import (
	"github.com/joho/godotenv"
)

//goland:noinspection GoUnhandledErrorResult
func init() {
	godotenv.Load()
}
