package main

import (
	"./nlb"
	"github.com/alexsasharegan/dotenv"
	"log"
)

func main() {
	err := dotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	bookTitle := "Kiasunomics©:Stories of Singaporean Economic Behaviours"
	availability := nlb.GetAvailabilityByTitle(bookTitle)
	log.Println(availability)
}
