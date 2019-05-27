package nlb

import (
	"github.com/tiaguinho/gosoap"
	"log"
)

const nlbURL = "http://openweb-stg.nlb.gov.sg/OWS/CatalogueService.svc?wsdl"

var nlbAPIKey string

// Setup initialises the NLB WDSL client
func Setup(apiKey string) {
	nlbAPIKey = apiKey
}

// GetAvailability fetches availability of a given book
func GetAvailability(isbn string) {
	soap, err := gosoap.SoapClient(nlbURL)
	if err != nil {
		log.Fatalf("SoapClient error: %s", err)
	}

	soapParams := gosoap.Params{
		"APIKey": nlbAPIKey,
		"ISBN":   isbn,
	}

	var res *gosoap.Response

	res, err = soap.Call("GetIpLocation", soapParams)
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}

	log.Println(res)
}
