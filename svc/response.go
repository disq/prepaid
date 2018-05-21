package svc

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// ApiGW converts an (interface{},error) pair to an (APIGatewayProxyResponse,error) pair. All errors are converted to HTTP 500s and logged.
func (se *Service) ApiGW(data interface{}, err error) (events.APIGatewayProxyResponse, error) {
	if err != nil {
		se.logger.Printf("Returning error: %v", err)

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		}, nil
	}

	str, err := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		Body:       string(str) + "\n",
	}, err
}
