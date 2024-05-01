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
		if key == "objid" {
			temp_key := "oid"
			query_params = append(query_params, temp_key, value)
		} else {
			query_params = append(query_params, key, value)
		}
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
			fmt.Printf("ERROR: Query failed %s\n", err.Error())
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
		}
		defer rows.Close()

		var audits []QueriedAudit
		for rows.Next() {
			var currentAudit QueriedAudit
			if err := rows.Scan(&currentAudit.Who, &currentAudit.Oid, &currentAudit.Namespace, &currentAudit.FieldName,
				&currentAudit.Before, &currentAudit.After, &currentAudit.EventTime); err != nil {
				fmt.Printf("ERROR: rows.Scan() failed %s\n", err.Error())
				return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
			}
			audits = append(audits, currentAudit)
		}
		b_response, err := json.Marshal(audits)
		if err != nil {
			fmt.Printf("ERROR: json.Marshal() failed %s\n", err.Error())
		}
		s_response := string(b_response)
		return events.APIGatewayProxyResponse{Body: s_response, StatusCode: 200}, nil
	}

	if query_params[0] == "namespace" && query_params[2] == "oid" {
		rows, err := db.Query("SELECT who, oid, namespace, field_name, before, after, event_time FROM audits where namespace = $1 and oid = $2 ORDER BY event_time desc",
			query_params[1], query_params[3])
		if err != nil {
			fmt.Printf("ERROR: Query failed %s\n", err.Error())
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
		}
		defer rows.Close()

		var audits []QueriedAudit
		for rows.Next() {
			var currentAudit QueriedAudit
			if err := rows.Scan(&currentAudit.Who, &currentAudit.Oid, &currentAudit.Namespace, &currentAudit.FieldName,
				&currentAudit.Before, &currentAudit.After, &currentAudit.EventTime); err != nil {
				fmt.Printf("ERROR: rows.Scan() failed %s\n", err.Error())
				return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, err
			}
			audits = append(audits, currentAudit)
		}
		b_response, err := json.Marshal(audits)
		if err != nil {
			fmt.Printf("ERROR: json.Marshal() failed %s\n", err.Error())
		}
		s_response := string(b_response)
		return events.APIGatewayProxyResponse{Body: s_response, StatusCode: 200}, nil
	}

	return events.APIGatewayProxyResponse{Body: http.StatusText(400), StatusCode: 400}, nil
}

//
// end of file
//
