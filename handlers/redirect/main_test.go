package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/yaizuuuu/url-shortener-lambda-go/db"
	"net/http"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		path, method string
		status       int
	}{
		{"xKlNKGomg", http.MethodGet, http.StatusPermanentRedirect},
		{"xKlNKGomg", http.MethodPost, http.StatusBadRequest},
		{"xKlNKGomg", http.MethodGet, http.StatusNotFound},
	}

	for _, te := range tests {
		res, _ := handler(events.APIGatewayProxyRequest{
			PathParameters: map[string]string{"shorten_resource": te.path},
			HTTPMethod:     te.method,
		})

		if res.StatusCode != te.status {
			t.Errorf("ExitStatus=%d, want %d", res.StatusCode, te.status)
		}
	}
}

type Link struct {
	ShortenResource string `json:"shorte_resource"`
	OriginalURL      string `json:"original_url"`
}

func prepare()  {
	DynamoDB = db.TestNew()

	if err := DynamoDB.CreateLinkTable(); err != nil {
		panic(err)
	}

	link := &Link{
		ShortenResource: "xKlNKGong",
		OriginalURL: "https://example.com",
	}
	_, err := DynamoDB.PutItem(link)
	if err != nil {
		panic(err)
	}
}

func cleanUp() {
	if err := DynamoDB.DeleteLinkTable(); err != nil {
		panic(err)
	}
	DynamoDB = db.DB{}
}

func TestMain(m *testing.M)  {
	prepare()
	exitCode := m.Run()
	cleanUp()
	os.Exit(exitCode)
}
