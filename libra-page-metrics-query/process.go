//
// main message processing
//

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	_ "github.com/lib/pq"
	"net/http"
)

type ObjectMetrics struct {
	Oid       string        `json:"oid"`
	Namespace string        `json:"namespace"`
	ViewCount int           `json:"views"`
	Files     []BlobMetrics `json:"files,omitempty "`
}

type BlobMetrics struct {
	TargetId      string `json:"target_id"`
	DownloadCount int    `json:"downloads"`
}

func process(messageId string, messageSrc string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var namespace string
	var oid string

	// log inbound query parameters
	for key, value := range request.QueryStringParameters {
		fmt.Printf("DEBUG: query param [%s] = [%s]\n", key, value)
		switch key {
		case "namespace":
			namespace = value
		case "oid":
			oid = value
		}
	}

	// log inbound headers
	for key, value := range request.Headers {
		fmt.Printf("DEBUG: header [%s] = [%s]\n", key, value)
	}

	// ensure we have the parameters we need
	if len(namespace) == 0 || len(oid) == 0 {
		err := fmt.Errorf("Missing required query params: [namespace, oid]")
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusBadRequest}, err
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

	// get a count of view events for this object
	var viewCount int
	// select count(*) from page_metrics where namespace = ns and oid = oid and metric_type = 'view'
	err = db.QueryRow("SELECT COUNT(*) FROM page_metrics WHERE namespace = $1 and oid = $2 and metric_type = 'view'", namespace, oid).Scan(&viewCount)
	if err != nil {
		fmt.Printf("ERROR: query failed (%s)\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}

	// get a list of filenames and the corresponding count of their download events
	var blobMetrics []BlobMetrics
	// select select target_id, count(*) from page_metrics where namespace = ns and oid = oid and metric_type = 'download' group by 1
	rows, err := db.Query("SELECT target_id, COUNT(*) FROM page_metrics WHERE namespace = $1 and oid = $2 and metric_type = 'download' GROUP BY 1", namespace, oid)
	if err != nil {
		fmt.Printf("ERROR: query failed (%s)\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}
	defer rows.Close()

	for rows.Next() {
		var fileMetrics BlobMetrics
		if err := rows.Scan(&fileMetrics.TargetId, &fileMetrics.DownloadCount); err != nil {
			fmt.Printf("ERROR: rows.Scan() failed (%s)\n", err.Error())
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
		}
		blobMetrics = append(blobMetrics, fileMetrics)
	}

	// check to see if we have results
	if viewCount == 0 && len(blobMetrics) == 0 {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
	}

	// construct the response
	response := ObjectMetrics{}
	response.Namespace = namespace
	response.Oid = oid
	response.ViewCount = viewCount
	response.Files = blobMetrics

	buf, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("ERROR: json.Marshal() failed (%s)\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}
	fmt.Printf("DEBUG: response [%s]\n", string(buf))
	return events.APIGatewayProxyResponse{Body: string(buf), StatusCode: http.StatusOK}, nil
}

//
// end of file
//
