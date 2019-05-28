package routes

import (
	"../goodreads"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"net/url"
	"os"
)

// GoodreadsAuthorise directs the user to gr authentication
func GoodreadsAuthorise(c echo.Context) error {
	apiKey := os.Getenv("GOODREADS_API_KEY")
	apiSecret := os.Getenv("GOODREADS_API_SECRET")
	callback := os.Getenv("GOODREADS_CALLBACK")

	oauthClient := goodreads.NewClient(apiKey, apiSecret)
	tempCred, err := oauthClient.RequestTemporaryCredentials(nil, "https://bbb2f57b.ngrok.io/goodreads_callback", nil)
	if err != nil {
		return c.String(http.StatusBadRequest, "Error requesting token "+err.Error())
	}
	log.Println(tempCred)

	redirectURL := oauthClient.AuthorizationURL(tempCred, url.Values{"oauth_callback": {callback}})
	return c.Redirect(http.StatusFound, redirectURL)
}

// GoodreadsCallback handles the gr callback
func GoodreadsCallback(c echo.Context) error {
	oauthToken := c.QueryParam("oauth_token")
	isAuthorised := c.QueryParam("authorize")

	log.Println(c.Request())

	if isAuthorised != "1" {
		return c.String(http.StatusUnauthorized, "Must allow goodlib to access Goodreads data")
	}

	availabilityURL, _ := url.Parse(c.Scheme() + "://" + c.Request().Host + "/availability")

	q := availabilityURL.Query()
	q.Set("oauth_token", oauthToken)

	return c.Redirect(http.StatusFound, availabilityURL.String())
}
