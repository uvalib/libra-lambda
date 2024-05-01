//
// main message processing
//

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/lib/pq"
)

type QueriedAudit struct {
	Who       *string    `json:"who"`
	Oid       *string    `json:"oid"`
	Namespace *string    `json:"namespace"`
	FieldName *string    `json:"fieldName"`
	Before    *string    `json:"before"`
	After     *string    `json:"after"`
	EventTime *time.Time `json:"eventTime"`
}

func process(messageId string, messageSrc string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	query_params := make([]string, 0, len(request.QueryStringParameters))
	// log inbound query parameters
	for key, value := range request.QueryStringParameters {
		fmt.Printf("DEBUG: query param [%s] = [%s]\n", key, value)
		if key == "objid" {
			temp_key := "oid"
			query_params = append(query_params, temp_key, value)
		} else {
			query_params = append(query_params, key, value)
		}
	}

	// log inbound headers
	for key, value := range request.Headers {
		fmt.Printf("DEBUG: header [%s] = [%s]\n", key, value)
	}

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}

	connectionStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		fmt.Printf("ERROR: unable to open database (%s)\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}

	// cleanup
	defer db.Close()

	// Adjust work query if the args are reversed from expected order
	if query_params[0] == "oid" && query_params[2] == "namespace" {
		query_params[0] = "namespace"
		query_params[2] = "oid"
		temp := query_params[1]
		query_params[1] = query_params[3]
		query_params[3] = temp
	}

	if query_params[0] == "who" {
		rows, err := db.Query("SELECT who, oid, namespace, field_name, before, after, event_time FROM audits where who = $1 ORDER BY event_time desc",
			query_params[1])
		if err != nil {
			fmt.Printf("ERROR: query failed (%s)\n", err.Error())
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
		}
		defer rows.Close()

		var audits []QueriedAudit
		for rows.Next() {
			var currentAudit QueriedAudit
			if err := rows.Scan(&currentAudit.Who, &currentAudit.Oid, &currentAudit.Namespace, &currentAudit.FieldName,
				&currentAudit.Before, &currentAudit.After, &currentAudit.EventTime); err != nil {
				fmt.Printf("ERROR: rows.Scan() failed (%s)\n", err.Error())
				return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
			}
			audits = append(audits, currentAudit)
		}
		b_response, err := json.Marshal(audits)
		if err != nil {
			fmt.Printf("ERROR: json.Marshal() failed (%s)\n", err.Error())
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
		}
		status := http.StatusOK
		if len(audits) == 0 {
			status = http.StatusNotFound
		}
		fmt.Printf("INFO: returning %d row(s)\n", len(audits))
		return events.APIGatewayProxyResponse{Body: string(b_response), StatusCode: status}, nil
	}

	if query_params[0] == "namespace" && query_params[2] == "oid" {
		rows, err := db.Query("SELECT who, oid, namespace, field_name, before, after, event_time FROM audits where namespace = $1 and oid = $2 ORDER BY event_time desc",
			query_params[1], query_params[3])
		if err != nil {
			fmt.Printf("ERROR: query failed (%s)\n", err.Error())
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
		}
		defer rows.Close()

		var audits []QueriedAudit
		for rows.Next() {
			var currentAudit QueriedAudit
			if err := rows.Scan(&currentAudit.Who, &currentAudit.Oid, &currentAudit.Namespace, &currentAudit.FieldName,
				&currentAudit.Before, &currentAudit.After, &currentAudit.EventTime); err != nil {
				fmt.Printf("ERROR: rows.Scan() failed (%s)\n", err.Error())
				return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
			}
			audits = append(audits, currentAudit)
		}
		b_response, err := json.Marshal(audits)
		if err != nil {
			fmt.Printf("ERROR: json.Marshal() failed (%s)\n", err.Error())
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
		}
		status := http.StatusOK
		if len(audits) == 0 {
			status = http.StatusNotFound
		}
		fmt.Printf("INFO: returning %d row(s)\n", len(audits))
		return events.APIGatewayProxyResponse{Body: string(b_response), StatusCode: status}, nil
	}

	return events.APIGatewayProxyResponse{Body: http.StatusText(http.StatusBadRequest), StatusCode: http.StatusBadRequest}, nil
}

//
// end of file
//
