//
// main message processing
//

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
)

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	// convert to librabus event
	ev, err := uvalibrabus.MakeBusEvent(rawMsg)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling libra bus event (%s)\n", err.Error())
		return err
	}

	fmt.Printf("INFO: EVENT %s from:%s -> %s\n", messageId, messageSrc, ev.String())

	audit, err := uvalibrabus.MakeAuditEvent(ev.Detail)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling libra audit event (%s)\n", err.Error())
		return err
	}

	fmt.Printf("INFO: Audit %v\n", audit)

	db, err := sql.Open("postgres", os.Getenv("DB_CONNECTION"))
	if err != nil {
		panic(err)
	}

	result, err := db.Exec("INSERT INTO audits (who, oid, namespace, field_name, before, after) VALUES ($1, $2, $3, $4, $5, $6)",
		audit.Who,
		ev.Identifier,
		ev.Namespace,
		audit.FieldName,
		audit.Before,
		audit.After,
	)

	if err != nil {
		panic(fmt.Errorf("ERROR: %s", err))
	}

	n, err := result.RowsAffected()
	if err != nil {
		panic(fmt.Errorf("ERROR: %s", err))
	}

	fmt.Printf("INFO: Inserted %d row\n", n)
	return nil
}

//
// end of file
//
