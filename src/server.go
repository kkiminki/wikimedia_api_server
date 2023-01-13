package main

import (
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// Static URL for the Wikimedia API
const WikimediaUrl = "https://en.wikipedia.org/w/api.php?"


// Start the server and setup the routes
func main() {
	router := gin.Default()
	router.GET("/person", getPerson)

	router.Run("localhost:8080")
}

// Make a call to the Wikimedia API to retrieve the
// data for the passed in name
func retrievePersonData(name string) (string){

	// Coble together the request string
	base, err := url.Parse(WikimediaUrl)
    if err != nil {
        return "blah"
    }

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

	base.RawQuery = p.Encode()

	// Initialize http client and execute the get request
	c := http.Client{}
	resp, err := c.Get(base.String())
    if err != nil {
        fmt.Printf("Error %s", err)
        return "blah"
    }
	
	// Parse out the response body and return it
	body, err := ioutil.ReadAll(resp.Body)
    return string(body[:])

}

// Route handler to get the short description
// on the person passed in through the request
// url
func getPerson(c *gin.Context) {
	name := c.Query("name")

	contents := retrievePersonData(name)

	c.String(http.StatusOK, "%s", contents)
}
