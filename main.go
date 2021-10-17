// main.go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

var isLambda bool

func init() {
	isLambda = len(os.Getenv("_LAMBDA_SERVER_PORT")) > 0
	log.SetReportCaller(true)
	if isLambda {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
}

func hello(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lp := request.QueryStringParameters["lp"]

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello Amazonian World %s!", lp),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(hello)
}
