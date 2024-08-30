package dcmenu_test

import (
	"encoding/json"
	"testing"

	"github.com/jackmerrill/hampbot/internal/utils/dcmenu"
)

func TestMain(t *testing.T) {
	t.Log("Testing Main")
}

func TestParseWebsite(t *testing.T) {
	t.Log("Testing ParseWebsite")

	url, err := dcmenu.ParseWebsite()

	if err != nil {
		t.Error(err)
	}

	if url == nil {
		t.Error("url is nil")
	}

	if *url == "" {
		t.Error("url is empty")
	}

	t.Log(*url)
}

func TestParseURL(t *testing.T) {
	t.Log("Testing ParseURL")

	url, err := dcmenu.ParseWebsite()

	if err != nil {
		t.Error(err)
	}

	if url == nil {
		t.Error("url is nil")
	}

	if *url == "" {
		t.Error("url is empty")
	}

	parsedURL, err := dcmenu.ParseURL(*url)

	if err != nil {
		t.Error(err)
	}

	if parsedURL == nil {
		t.Error("parsedURL is nil")
	}

	if *parsedURL == "" {
		t.Error("parsedURL is empty")
	}

	t.Log(*parsedURL)
}

func TestParseCSV(t *testing.T) {
	t.Log("Testing ParseCSV")

	url, err := dcmenu.ParseWebsite()

	if err != nil {
		t.Error(err)
	}

	if url == nil {
		t.Error("url is nil")
	}

	if *url == "" {
		t.Error("url is empty")
	}

	parsedURL, err := dcmenu.ParseURL(*url)

	if err != nil {
		t.Error(err)
	}

	if parsedURL == nil {
		t.Error("parsedURL is nil")
	}

	if *parsedURL == "" {
		t.Error("parsedURL is empty")
	}

	csv, err := dcmenu.ParseCSV(*parsedURL)

	if err != nil {
		t.Error(err)
	}

	if csv == nil {
		t.Error("csv is nil")
	}

	if len(*csv) == 0 {
		t.Error("csv is empty")
	}

	json, err := json.MarshalIndent(csv, "", "  ")

	if err != nil {
		t.Error(err)
	}

	t.Log(string(json))
}
