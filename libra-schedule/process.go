//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
)

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	fmt.Printf("EVENT %s from:%s -> %s\n", messageId, messageSrc, string(rawMsg))
	return nil
}

//
// end of file
//
