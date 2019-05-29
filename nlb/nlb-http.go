package nlb

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"net/url"
	s "strings"
)

const nlbRootURL = "https://catalogue.nlb.gov.sg"
const nlbSearchCatalogue = "https://catalogue.nlb.gov.sg/cgi-bin/spydus.exe/ENQ/WPAC/BIBENQ"

// AvailabilityInfo models the information about the status of the item at each branch
type AvailabilityInfo struct {
	BranchCode string
	BranchName string
	Status     string
}

// Book models the info returned from the NLB API
type Book struct {
	Title        string
	CallNumber   string
	Availability []AvailabilityInfo
}

// GetBookURLByTitle searches NLB's catalogue by title and returns the URL of the book
// ENTRY_NAME=BS is basic search, ENTRY_NAME=TI is title search
// ENTRY_TYPE=K is keyword, ENTRY_TYPE=E is exact search
// QRYTEXT specifies category of items
func GetBookURLByTitle(title string, isExact bool) (string, error) {
	entryType := "K"
	if isExact {
		entryType = "E"
	}
	resp, _, errs := gorequest.New().
		Get(nlbSearchCatalogue).
		Query("ENTRY_NAME=BS&ENTRY_TYPE=" + entryType + "&QRYTEXT=Books").
		Query("ENTRY=" + url.QueryEscape(title)).
		End()
	if errs != nil {
		return "", errors.New("Error searching catalogue: " + errs[0].Error())
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Catalogue returned invalid status code: " + resp.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", errors.New("Invalid catalogue document: " + err.Error())
	}

	matchingEls := doc.Find("#result-content-grid .card-body .card-text.availability .col-12 a")
	bookURL, exists := matchingEls.First().Attr("href")

	if !exists {
		return "", errors.New("Book not found")
	}
	return nlbRootURL + bookURL, nil
}

// GetAvailabilityByURL takes a book URL and returns availability info
func GetAvailabilityByURL(url string) (Book, error) {
	transformedURL := s.Replace(url, "?RECDISP=REC", "", 1)
	resp, _, errs := gorequest.New().
		Get(transformedURL).
		End()
	if errs != nil {
		return Book{}, errors.New("Error fetching availability: " + errs[0].Error())
	}
	if resp.StatusCode != http.StatusOK {
		return Book{}, errors.New("Availability URL returned invalid status code: " + resp.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return Book{}, errors.New("Invalid availability page: " + err.Error())
	}

	data := Book{}
	doc.Find("table tbody tr").Each(func(i int, el *goquery.Selection) {
		if data.CallNumber == "" {
			spans := el.Find("td").Eq(2).Find("span")
			data.CallNumber = spans.Eq(0).Text() + " " + spans.Eq(1).Text()
		}
		branch, _ := el.Find("td").Eq(1).Find("book-location").Attr("data-branch")
		availability := el.Find("td").Eq(3).Text()
		data.Availability = append(data.Availability, AvailabilityInfo{
			BranchCode: branch,
			Status:     availability,
		})
	})
	return data, nil
}

// GetAvailabilityByTitle takes a book title and returns the availability info
func GetAvailabilityByTitle(title string, isExact bool) (Book, error) {
	url, getURLErr := GetBookURLByTitle(title, isExact)
	if getURLErr != nil {
		return Book{}, getURLErr
	}
	book, getBookErr := GetAvailabilityByURL(url)
	if getBookErr != nil {
		return Book{}, getBookErr
	}
	return book, nil
}
