package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"errors"
	"bytes"

	"github.com/gin-gonic/gin"
)

// Static URL for the Wikimedia API
const WikimediaUrl = "https://en.wikipedia.org/w/api.php?"

type APIClient struct {
	URL			string
	client		*http.Client
}

func (c *APIClient) Get(params url.Values) (*http.Response, error) {

	// Coble together the request string
	base, err := url.Parse(c.URL)
	if err != nil {

		// If it bails out throw an internal server error response
		t := http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("Internal Server Error")),
		}
		return &t, errors.New(fmt.Sprintf("Failed to parse url: %s", c.URL))
	}

	base.RawQuery = params.Encode()
	return c.client.Get(base.String())

}

// Start the server and setup the routes
func main() {
	router := gin.Default()
	router.GET("/person", getPerson)

	router.Run("localhost:8080")
}

// Make a call to the Wikimedia API to retrieve the
// data for the passed in name
func retrievePersonData(name string, client *APIClient) (int, string) {



	// Add on parameters for the request, the name gets
	// passed in here
	p := url.Values{}
	p.Add("action", "query")
	p.Add("prop", "revisions")
	p.Add("titles", name)
	p.Add("rvlimit", "1")
	p.Add("formatversion", "2")
	p.Add("format", "json")
	p.Add("rvprop", "content")


	// Initialize http client and execute the get request
	resp, err := client.Get(p)
	if err != nil {
		fmt.Printf("Error %s", err)
		return http.StatusInternalServerError, fmt.Sprintf("Error attempting to get data for %s from %s", name, WikimediaUrl)
	}

	// Parse out the response body and return it
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error pasing response %s", err)
		return http.StatusInternalServerError, fmt.Sprintf("Error attempting to get data for %s from %s", name, WikimediaUrl)	
	}

	defer resp.Body.Close()
	return http.StatusOK, string(body[:])

}

// Parse the data from Wikimedia and return the
// persons name and a short description
func parsePersonData(name string, data string) string {

	// Search for a Short description block in the wikimedia data
	re := regexp.MustCompile(`\{\{(Short description)\|(.*?)\}\}`)
	match := re.FindStringSubmatch(data)

	// If there was a match for the wildcard return it
	// otherwise report back that nothing was found
	if len(match) > 2 {
		return match[2]
	} else {
		return fmt.Sprintf("Description for %s could not be found", name)
	}

}

// Route handler to get the short description
// on the person passed in through the request
// url
func getPerson(c *gin.Context) {
	name := c.Query("name")

	client := APIClient {WikimediaUrl, &http.Client{}}
	status, contents := retrievePersonData(name, &client)
	if status != http.StatusOK {
		c.String(status, contents)
	} else {
		description := parsePersonData(name, contents)
		c.String(status, "{\"%s\": \"%s\"}", name, description)
	}
}
