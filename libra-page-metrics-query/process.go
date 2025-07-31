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

type ObjectMetrics struct {
	Oid       string        `json:"oid"`
	Namespace string        `json:"namespace"`
	ViewCount int           `json:"views"`
	Files     []BlobMetrics `json:"files,omitempty "`
}

type BlobMetrics struct {
	FileId        string `json:"file_id"`
	DownloadCount int    `json:"downloads"`
}

type DbQueryResponse struct {
	MetricType string
	Oid        string
	Namespace  string
	FileId     string
	SourceIp   string
	UserAgent  string
	AcceptLang string
	EventTime  time.Time
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

	// select count(*) from page_metrics where namespace = ns, oid = oid, mtype = 'views'

	// select select file_id, count(*) from page_metrics where namespace = ns, oid = oid, mtype = 'downloads' group by 1

	//var metrics []DbQueryResponse
	//		rows, err := db.Query("SELECT who, oid, namespace, field_name, before, after, event_time FROM audits where namespace = $1 and oid = $2 ORDER BY event_time desc",
	//			query_params[1], query_params[3])
	//		if err != nil {
	//			fmt.Printf("ERROR: query failed (%s)\n", err.Error())
	//			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	//		}
	//		defer rows.Close()
	//
	//		var audits []QueriedAudit
	//		for rows.Next() {
	//			var currentAudit QueriedAudit
	//			if err := rows.Scan(&currentAudit.Who, &currentAudit.Oid, &currentAudit.Namespace, &currentAudit.FieldName,
	//				&currentAudit.Before, &currentAudit.After, &currentAudit.EventTime); err != nil {
	//				fmt.Printf("ERROR: rows.Scan() failed (%s)\n", err.Error())
	//				return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	//			}
	//			audits = append(audits, currentAudit)
	//		}

	// check to see if we have results
	//if len(metrics) == 0 {
	//	return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound}, nil
	//}

	response := ObjectMetrics{}
	response.Namespace = namespace
	response.Oid = oid
	response.ViewCount = 99
	bm := BlobMetrics{FileId: "this", DownloadCount: 5}
	response.Files = []BlobMetrics{bm}

	buf, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("ERROR: json.Marshal() failed (%s)\n", err.Error())
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, err
	}
	return events.APIGatewayProxyResponse{Body: string(buf), StatusCode: http.StatusOK}, nil
}

//
// end of file
//
