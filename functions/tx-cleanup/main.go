package main

import (
	"context"
	"log"
	"os"

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
func Handler() error {
	return se.TxCleanup()
}
