package main

import (
	"github.com/joho/godotenv"
	"tuxedo-email-service/services"
)

func main() {
	godotenv.Load()
	services.Stream()
}
