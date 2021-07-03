package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/yaizuuuu/url-shortener-lambda-go/db"
	"net/http"
)

var DynamoDB db.DB

func init() {
	DynamoDB = db.New()
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r, err := parseRequest(request)
	if err != nil {
		return response(
			http.StatusBadRequest,
			errorResponseBody(err.Error()),
		), nil
	}

	URL, err := DynamoDB.GeiItem(r)
	if err != nil {
		return response(
			http.StatusInternalServerError,
			errorResponseBody(err.Error()),
		), nil
	}
	if URL == "" {
		return response(http.StatusNotFound, ""), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusPermanentRedirect,
		Headers: map[string]string{
			"location": URL,
		},
	}, nil
}

func parseRequest(req events.APIGatewayProxyRequest) (string, error) {
	if req.HTTPMethod != http.MethodGet {
		return "", fmt.Errorf("use GET request")
	}

	shorternResource := req.PathParameters["shorten_resource"]

	return shorternResource, nil
}

func response(code int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

func errorResponseBody(msg string) string {
	return fmt.Sprintf("{\"message\": \"%s\"}", msg)
}
