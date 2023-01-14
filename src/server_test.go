package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

// This test checks checkIfMissing against
// valid input
func TestCheckIfMissingValid(t *testing.T) {
	data := `{"needs": "validJson", "pages": [{"missing": "false"}]}`
	status := checkIfMissing(data)
	if status != http.StatusOK {
		t.Fatalf("checkIfMissing returned unexpected value for valid input: %d", status)
	}
}

// This test checks checkIfMissing against
// valid input
func TestCheckIfMissingInvalid(t *testing.T) {
	data := `{"needs": "validJson", "pages": [{"missing": "true"}]}`
	status := checkIfMissing(data)
	if status != http.StatusBadRequest {
		t.Fatalf("checkIfMissing returned unexpected value for invalid input: %d", status)
	}
}

// This test checks normalizeName against
// underscore separated input
func TestNormalizeUnderscoreInput(t *testing.T) {
	input := "test_name"
	normalizedInput := normalizeName(input, "_")
	if normalizedInput != "Test_Name" {
		t.Fatalf("normalizeName returned unexpected value for underscore separated input: %s", normalizedInput)
	}
}

// This test checks normalizeName against
// space separated input
func TestNormalizeSpaceInput(t *testing.T) {
	input := "test name"
	normalizedInput := normalizeName(input, " ")
	if normalizedInput != "Test Name" {
		t.Fatalf("normalizeName returned unexpected value for underscore separated input: %s", normalizedInput)
	}
}

// This test checks normalizeName against
// input that doesn't need to be normalized
func TestNormalizeInput(t *testing.T) {
	input := "Testname"
	normalizedInput := normalizeName(input, " ")
	if normalizedInput != "Testname" {
		t.Fatalf("normalizeName returned unexpected value for underscore separated input: %s", normalizedInput)
	}
}
