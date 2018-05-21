package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
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

	amt := float64(0)

	if amtStr := strings.TrimSpace(request.QueryStringParameters["amt"]); amtStr != "" {
		res, err := strconv.ParseFloat(amtStr, 64)
		if err != nil {
			return se.ApiGW(nil, fmt.Errorf("Invalid amt"))
		}
		amt = res
	} else {
		return se.ApiGW(nil, fmt.Errorf("Invalid amt"))
	}

	merchant := strings.TrimSpace(request.QueryStringParameters["merchant"])
	if merchant == "" {
		return se.ApiGW(nil, fmt.Errorf("Invalid merchant"))
	}

	return se.ApiGW(se.NewTx(id, amt, merchant))
}
