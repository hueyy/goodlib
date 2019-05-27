package main

import (
	"./nlb"
	"github.com/alexsasharegan/dotenv"
	"log"
	"os"
)

func main() {
	err := dotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	nlb.Setup(os.Getenv("NLB_API_KEY"))

	nlb.GetAvailability("0062433652")
}
