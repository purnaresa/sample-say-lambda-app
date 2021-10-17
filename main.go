// main.go
package main

import (
	"context"
	"fmt"
	"os"

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

type MyEvent struct {
	Name string `json:"name"`
}

func hello(ctx context.Context, name MyEvent) (string, error) {
	return fmt.Sprintf("Hello Amazonian World %s!", name.Name), nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	if isLambda == true {
		lambda.Start(hello)

	} else {
		out, err := hello(nil, MyEvent{Name: "a"})
		if err != nil {
			log.Fatalln(err)
		}
		log.Infoln(out)
	}
}
