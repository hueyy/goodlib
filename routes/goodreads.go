package routes

import (
	"../goodreads"
	"github.com/labstack/echo"
	"net/http"
	"os"
)

// GoodreadsAuthorise directs the user to gr authentication
func GoodreadsAuthorise(c echo.Context) error {
	apiKey := os.Getenv("GOODREADS_API_KEY")
	apiSecret := os.Getenv("GOODREADS_API_SECRET")
	callback := os.Getenv("BASE_URL") + "/goodreads_callback"

	redirectURL, err := goodreads.GetAuthenticationURL(apiKey, apiSecret, callback)
	if err != nil {
		return c.String(http.StatusBadRequest, "Error requesting temporary credentials "+err.Error())
	}

	return c.Redirect(http.StatusFound, redirectURL)
}

// GoodreadsCallback handles the gr callback
func GoodreadsCallback(c echo.Context) error {
	oauthToken := c.QueryParam("oauth_token")
	isAuthorised := c.QueryParam("authorize")

	if isAuthorised != "1" {
		return c.String(http.StatusUnauthorized, "Must allow goodlib to access Goodreads data")
	}

	token, err := goodreads.GetTokenCredentials(oauthToken)
	if err != nil {
		return c.String(http.StatusBadRequest, "Error requesting token")
	}

	availabilityURL := c.Scheme() + "://" + c.Request().Host + "/availability?token=" + token

	return c.Redirect(http.StatusFound, availabilityURL)
}
