package main

import (
	// "./nlb"
	"./routes"
	"github.com/alexsasharegan/dotenv"
	"github.com/hueyy/go-reads"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"os"
)

func getToReadBooks() {
	goodreadsAPIKey := os.Getenv("GOODREADS_API_KEY")
	grClient := goodreads.NewClient(goodreadsAPIKey)
	uid := "47976050"
	reviews, err := grClient.ReviewList(uid, "to-read", "date_added", "", "d", 1, 10)
	if err != nil {
		log.Println("Cannot fetch books on shelf: ")
		panic(err)
	}
	log.Println("Reviews:")
	for i, rev := range reviews {
		log.Printf(" %d. [%d stars, %s] %s\n", i+1, rev.Rating, rev.ReadAt, rev.Book.Title)
	}
}

func main() {
	err := dotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/availability", routes.ShowAvailability)
	e.GET("/goodreads_authorise", routes.GoodreadsAuthorise)
	e.GET("/goodreads_callback", routes.GoodreadsCallback)
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
	// getToReadBooks()
	// bookTitle := "Kiasunomics©:Stories of Singaporean Economic Behaviours"
	// availability := nlb.GetAvailabilityByTitle(bookTitle)
	// log.Println(availability)
}
