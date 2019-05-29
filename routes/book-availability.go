package routes

import (
	"../goodreads"
	"../nlb"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

type bookInfo struct {
	Title        string
	NLBTitle     string
	Author       string
	CallNumber   string
	Availability []nlb.AvailabilityInfo
}

// ShowAvailability displays the availability of books in the user's to-read shelf
func ShowAvailability(c echo.Context) error {
	token := c.QueryParam("token")
	userID, err := goodreads.GetAuthenticatedUserID(token)
	if err != nil {
		return c.String(http.StatusBadRequest, "Could not get logged-in user")
	}
	books, err := goodreads.GetShelf(userID, token)
	if err != nil {
		return c.String(http.StatusBadRequest, "Could not fetch books on shelf")
	}

	type empty struct{}

	availabilitySlice := []bookInfo{}
	missing := []int{}
	queue := make(chan empty, len(books))
	for i, book := range books {
		availabilitySlice = append(availabilitySlice, bookInfo{
			Title:  book.Title,
			Author: book.Authors[0].Name,
		})
		go func(index int, bk goodreads.Book) {
			nlbData, err := nlb.GetAvailabilityByTitle(bk.Title, true)
			if err != nil {
				log.Println(bk.Title + " - " + err.Error())
				missing = append(missing, index)
			} else {
				availabilitySlice[index].NLBTitle = nlbData.Title
				availabilitySlice[index].CallNumber = nlbData.CallNumber
				availabilitySlice[index].Availability = nlbData.Availability
			}
			queue <- empty{}
		}(i, book)
	}
	for i := 0; i < len(availabilitySlice); i++ {
		<-queue
	}

	log.Println("retrying for " + string(len(missing)) + " books")
	success := 0
	for _, index := range missing {
		book := books[index]
		nlbData, err := nlb.GetAvailabilityByTitle(book.Title, false)
		if err != nil {
			log.Println(book.Title + " - " + err.Error())
			success++
		} else {
			availabilitySlice[index].NLBTitle = nlbData.Title
			availabilitySlice[index].CallNumber = nlbData.CallNumber
			availabilitySlice[index].Availability = nlbData.Availability
		}
	}
	log.Println(
		string(len(books)) + " total | " +
			string(len(missing)) + " not found | " +
			string(success) + " found later")
	return c.JSON(http.StatusOK, availabilitySlice)
}
