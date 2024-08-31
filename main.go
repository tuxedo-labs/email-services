package main

import (
	"fmt"
	"tuxedo-email-service/services"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Load .env")
	godotenv.Load()
	fmt.Println("Starting Service")
	services.Stream()
}
