package auth

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	log "github.com/victron/simpleLogger"
)

func RecieveCookies(baseUrl, username, password string) ([]*http.Cookie, error) {
	fullUrl, err := url.Parse(baseUrl)
	if err != nil {
		log.Error.Fatalln("Malformed URL: ", err)
	}
	values := url.Values{}
	values.Add("redirect", baseUrl)
	values.Add("username", username)
	values.Add("password", password)
	values.Add("submit", "Log+in+to+the+auction")
	fullUrl.RawQuery = values.Encode()

	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}

	res, err := client.PostForm(fullUrl.String(), values)
	if err != nil {
		log.Error.Fatalln(err)
	}
	defer res.Body.Close()

	for _, cookie := range cookieJar.Cookies(fullUrl) {
		log.Debug.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
	return cookieJar.Cookies(fullUrl), nil
}
