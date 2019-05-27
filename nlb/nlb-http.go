package nlb

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
	"net/url"
	s "strings"
)

const nlbRootURL = "https://catalogue.nlb.gov.sg"
const nlbSearchCatalogue = "https://catalogue.nlb.gov.sg/cgi-bin/spydus.exe/ENQ/WPAC/BIBENQ"

// GetBookURLByTitle searches NLB's catalogue by title and returns the URL of the book
func GetBookURLByTitle(title string) string {
	resp, _, errs := gorequest.New().
		Get(nlbSearchCatalogue).
		Query("ENTRY_NAME=BS&ENTRY_TYPE=K&QRYTEXT=Books").
		Query("ENTRY=" + url.QueryEscape(title)).
		End()
	if errs != nil {
		log.Println("Error searching catalogue", errs)
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("Catalogue returned invalid status code", resp.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal("Invalid catalogue document", err)
	}

	matchingEls := doc.Find("#result-content-grid .card-body .card-text.availability .col-12 a")
	bookURL, exists := matchingEls.First().Attr("href")

	if !exists {
		log.Fatal("Book not found", matchingEls)
	}
	return nlbRootURL + bookURL
}

// GetAvailabilityByURL takes a book URL and returns availability info
func GetAvailabilityByURL(url string) map[string]string {
	transformedURL := s.Replace(url, "?RECDISP=REC", "", 1)
	resp, _, errs := gorequest.New().
		Get(transformedURL).
		End()
	if errs != nil {
		log.Println("Error fetching availability", errs)
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("Availability URL returned invalid status code", resp.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal("Invalid availability page", err)
	}

	data := make(map[string]string)
	data["callNumber"] = ""
	doc.Find("table tbody tr").Each(func(i int, el *goquery.Selection) {
		if data["callNumber"] == "" {
			spans := el.Find("td").Eq(2).Find("span")
			data["callNumber"] = spans.Eq(0).Text() + " " + spans.Eq(1).Text()
		}
		branch, _ := el.Find("td").Eq(1).Find("book-location").Attr("data-branch")
		availability := el.Find("td").Eq(3).Text()
		data[branch] = availability
	})
	return data
}

// GetAvailabilityByTitle takes a book title and returns the availability info
func GetAvailabilityByTitle(title string) map[string]string {
	url := GetBookURLByTitle(title)
	return GetAvailabilityByURL(url)
}
