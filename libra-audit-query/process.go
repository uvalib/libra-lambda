//
// main message processing
//

package main

import (
	"database/sql"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	_ "github.com/lib/pq"
)

func process(messageId string, messageSrc string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// log inbound query parameters
	for key, value := range request.QueryStringParameters {
		fmt.Printf("Query Param %s: %s\n", key, value)
	}

	// log inbound headers
	for key, value := range request.Headers {
		fmt.Printf("Header %s: %s\n", key, value)
	}

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	connectionStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		fmt.Printf("ERROR: unable to open database %s\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	// cleanup
	defer db.Close()

	// do stuff

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

//
// end of file
//
