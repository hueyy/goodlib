package routes

import (
	"github.com/labstack/echo"
	"net/http"
)

// ShowAvailability displays the availability of books in the user's to-read shelf
func ShowAvailability(c echo.Context) error {
	// oauthToken := c.QueryParam("oauth_token")
	return c.String(http.StatusOK, "yay")
}
