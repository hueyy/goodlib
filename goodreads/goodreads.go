package goodreads

import (
	"encoding/xml"
	"errors"
	"github.com/gomodule/oauth1/oauth"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Response models GR responses
type Response struct {
	User    User     `xml:"user"`
	Book    Book     `xml:"book"`
	Books   []Book   `xml:"books>book"`
	Reviews []Review `xml:"reviews>review"`
}

// User models a GR user object
type User struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
	Link string `xml:"link"`
}

// Shelf models a GR shelf object (different from Review)
type Shelf struct {
	ID        string `xml:"id"`
	BookCount string `xml:"book_count"`
	Name      string `xml:"name"`
}

// Author models a GR author object (contained in Books)
type Author struct {
	ID   string `xml:"id"`
	Name string `xml:"name"`
	Link string `xml:"link"`
}

// Book models a GR book object (contained in Shelves)
type Book struct {
	ID       string   `xml:"id"`
	Title    string   `xml:"title"`
	Link     string   `xml:"link"`
	ImageURL string   `xml:"image_url"`
	NumPages string   `xml:"num_pages"`
	Format   string   `xml:"format"`
	Authors  []Author `xml:"authors>author"`
	ISBN     string   `xml:"isbn"`
}

// Review models a GR review object
type Review struct {
	Book   Book   `xml:"book"`
	Rating int    `xml:"rating"`
	ReadAt string `xml:"read_at"`
	Link   string `xml:"link"`
}

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://www.goodreads.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://www.goodreads.com/oauth/authorize",
	TokenRequestURI:               "https://www.goodreads.com/oauth/access_token",
	SignatureMethod:               oauth.HMACSHA1,
}

var credentialStore = make(map[string]*oauth.Credentials)

// GetAuthenticationURL returns the gr authentication URL to redirect to
func GetAuthenticationURL(apiKey, apiSecret, callbackURL string) (string, error) {
	oauthClient.Credentials = oauth.Credentials{
		Token:  apiKey,
		Secret: apiSecret,
	}
	tempCred, err := oauthClient.RequestTemporaryCredentials(nil, callbackURL, nil)
	if err != nil {
		return "", errors.New("Error requesting temporary credentials")
	}
	credentialStore[tempCred.Token] = tempCred

	redirectURL := oauthClient.AuthorizationURL(tempCred, url.Values{"oauth_callback": {callbackURL}})
	return redirectURL, nil
}

// GetTokenCredentials requests and saves token credentials using the temporary oauth token
// returns token key
func GetTokenCredentials(oauthToken string) (string, error) {
	tempCred := credentialStore[oauthToken]
	tokenCred, _, err := oauthClient.RequestToken(nil, tempCred, "")
	if err != nil {
		return "", errors.New("Error requesting token")
	}
	delete(credentialStore, oauthToken)
	credentialStore[tokenCred.Token] = tokenCred
	return tokenCred.Token, nil
}

func parseResponse(response *http.Response) (Response, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Response{}, errors.New("Could not read GR response body")
	}
	response.Body.Close()

	var parsedResponse Response
	xml.Unmarshal(body, &parsedResponse)
	return parsedResponse, nil
}

func apiGet(token, uri string) (*http.Response, error) {
	creds := credentialStore[token]
	resp, err := oauthClient.Get(nil, creds, uri, url.Values{})
	if err != nil {
		emptyResponse := http.Response{}
		return &emptyResponse, errors.New("Error making authenticated request " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		emptyResponse := http.Response{}
		return &emptyResponse, errors.New("GR API returned invalid status code" + resp.Status)
	}
	return resp, nil
}

func getData(token, uri string) (Response, error) {
	resp, err := apiGet(token, uri)
	if err != nil {
		return Response{}, err
	}
	response, err := parseResponse(resp)
	if err != nil {
		return Response{}, err
	}
	return response, nil
}

// GetAuthenticatedUserID returns the userid of the oauthenticated user
func GetAuthenticatedUserID(token string) (string, error) {
	resp, err := getData(token, "https://www.goodreads.com/api/auth_user")
	if err != nil {
		return "", errors.New("Error fetching authenticated user ID")
	}
	return resp.User.ID, nil
}

// GetShelf fetches the authenticated user's shelf
func GetShelf(userid, token string) ([]Book, error) {
	resp, err := getData(
		token,
		"https://www.goodreads.com/review/list/"+userid+".xml", // ?v=2&shelf=to-read&sort=date_added&order=d
	)
	if err != nil {
		return []Book{}, errors.New("Error fetching books on shelf")
	}
	return resp.Books, nil
}
