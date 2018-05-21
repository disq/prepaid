package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/disq/prepaid/svc"
)

var (
	se *svc.Service
	lo *log.Logger
)

func main() {
	var err error

	lo = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

	se, err = svc.New(context.Background(), lo)
	if err != nil {
		lo.Fatal(err)
	}

	lambda.Start(Handler)
}

// Handler is the Lambda entrypoint.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := strings.TrimSpace(request.PathParameters["id"])

	if id == "" {
		return se.ApiGW(nil, fmt.Errorf("Invalid id"))
	}

	return se.ApiGW(se.TxStatus(id))
}
