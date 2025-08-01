//
// main message processing
//

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aquilax/truncate"

	_ "github.com/lib/pq"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
)

// keep the following in-sync with the schema defined in the migrations
var maxTargetIdSize = 64
var maxSourceIpSize = 32
var maxReferrerSize = 255
var maxUserAgentSize = 255
var maxAcceptLanguageSize = 32

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	// convert to librabus event
	ev, err := uvalibrabus.MakeBusEvent(rawMsg)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling libra bus event (%s)\n", err.Error())
		return err
	}

	fmt.Printf("INFO: EVENT %s from %s -> %s\n", messageId, messageSrc, ev.String())

	content, err := uvalibrabus.MakeContentEvent(ev.Detail)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling libra content event (%s)\n", err.Error())
		return err
	}

	fmt.Printf("INFO: Content %v\n", content)

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return err
	}

	connectionStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		fmt.Printf("ERROR: unable to open database %s\n", err.Error())
		return err
	}
	// cleanup
	defer db.Close()

	parsedEventTime, err := time.Parse(time.RFC3339, ev.EventTime)
	if err != nil {
		fmt.Printf("ERROR: unable to parse event time %s\n", err.Error())
		return err
	}

	var metricType string
	switch ev.EventName {
	case uvalibrabus.EventContentView:
		metricType = "view"

	case uvalibrabus.EventContentDownload:
		metricType = "download"

	default:
		fmt.Printf("INFO: uninteresting event, ignoring\n")
		return nil
	}

	// endure we do not exceed any character limits
	content.TargetId = truncate.Truncate(content.TargetId, maxTargetIdSize, "...", truncate.PositionEnd)
	content.SourceIp = truncate.Truncate(content.SourceIp, maxSourceIpSize, "...", truncate.PositionEnd)
	content.Referrer = truncate.Truncate(content.Referrer, maxReferrerSize, "...", truncate.PositionEnd)
	content.UserAgent = truncate.Truncate(content.UserAgent, maxUserAgentSize, "...", truncate.PositionEnd)
	content.AcceptLanguage = truncate.Truncate(content.AcceptLanguage, maxAcceptLanguageSize, "...", truncate.PositionEnd)

	// insert into the database
	result, err := db.Exec("INSERT INTO page_metrics (metric_type, namespace, oid, target_id, source_ip, referrer, user_agent, accept_lang, event_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		metricType,
		ev.Namespace,
		ev.Identifier,
		content.TargetId,
		content.SourceIp,
		content.Referrer,
		content.UserAgent,
		content.AcceptLanguage,
		parsedEventTime,
	)

	if err != nil {
		fmt.Printf("ERROR: db insert %s", err)
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("ERROR: rows affected %s", err)
		return err
	}

	fmt.Printf("INFO: Inserted %d row\n", n)

	// log the happy news
	fmt.Printf("INFO: EVENT %s from %s processed OK\n", messageId, messageSrc)
	return nil
}

//
// end of file
//
