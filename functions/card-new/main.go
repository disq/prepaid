package main

import (
	"context"
	"fmt"
	"log"
	"math"
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

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	amt := float64(0)

	if amtStr := strings.TrimSpace(request.QueryStringParameters["amt"]); amtStr != "" {
		res, err := strconv.ParseFloat(amtStr, 64)
		if err != nil || res < 0 || math.IsNaN(res) || math.IsInf(res, 0) {
			return se.ApiGW(nil, fmt.Errorf("Invalid amt"))
		}
		amt = res
	}

	return se.ApiGW(se.NewCard(amt))
}
