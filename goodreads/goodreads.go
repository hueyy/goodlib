package goodreads

import (
	"github.com/gomodule/oauth1/oauth"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
)

// NewClient returns a Goodreads oauth client
func NewClient(apiKey, apiSecret string) oauth.Client {
	var oauthClient = oauth.Client{
		Credentials: oauth.Credentials{
			Token:  apiKey,
			Secret: apiSecret,
		},
		TemporaryCredentialRequestURI: "https://www.goodreads.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://www.goodreads.com/oauth/authorize",
		TokenRequestURI:               "https://www.goodreads.com/oauth/access_token",
		SignatureMethod:               oauth.HMACSHA1,
	}
	return oauthClient
}

func getAuthenticatedUserID() string {
	req := gorequest.New().
		Get("https://www.goodreads.com/api/index")
	resp, body, errs := req.Query("KEY=&").End()
	if errs != nil {
		log.Println("Error searching catalogue", errs)
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("Catalogue returned invalid status code", resp.Status)
	}
	return body
}
