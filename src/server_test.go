package main

import (
	"testing"
    "net/http/httptest"
    "net/http"
    
)

// This test checks parsePersonData against
// valid input
func TestParsePersonDataValid(t *testing.T) {
	data := "SomeJUnkData{{Short description|This should return}}But none of this"
	description := parsePersonData("JohnDoe", data)
	if description != "This should return" {
		t.Fatalf("parsePersonData returned unexpected value for valid input: %s", description)
	}
}

// This test checks parsePersonData against
// invalid input
func TestParsePersonDataInvalid(t *testing.T) {
	data := "SomeJUnkData{{Shortdescription|This should return}}But none of this"
	description := parsePersonData("JohnDoe", data)
	if description != "Description for JohnDoe could not be found" {
		t.Fatalf("parsePersonData returned unexpected value for invalid input: %s", description)
	}
}

// This test checks retrievePersonData against
// invalid input
func TestRetrievePersonDataInvalid(t *testing.T) {

    // Create a mock server to raise an error
    handler := func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Length", "1")
    }
    ts := httptest.NewServer(http.HandlerFunc(handler))
    
    c := APIClient{ts.URL, &http.Client{}}
    defer ts.Close()

    status, contents := retrievePersonData("Yoshua_Bengio", &c)

    if status != http.StatusInternalServerError {
        t.Fatalf("retrievePersonData returned unexpected value for invalid input: %s", contents)
    }
}

// This test checks retrievePersonData against
// valid input
func TestRetrievePersonDataValid(t *testing.T) {

    c := APIClient{WikimediaUrl, &http.Client{}}

    status, contents := retrievePersonData("Yoshua_Bengio", &c)

    if status != http.StatusOK {
        t.Fatalf("retrievePersonData returned unexpected value for valid input: %s", contents)
    }
}