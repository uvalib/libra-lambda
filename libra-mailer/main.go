//
//
//

package main

import (
	"context"
	"fmt"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
)

func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {

	//if event == nil {
	//	return fmt.Errorf("received nil event")
	//}

    // loop through possible messages
    for _, message := range sqsEvent.Records {
       // convert to an eventbus event
       var mbEvent events.EventBridgeEvent
       err := json.Unmarshal([]byte(message.Body), &mbEvent)
       if err != nil {
          fmt.Printf("ERROR: unmarshaling event bridge event (%s)", err.Error())
          return err
       }

       // convert to librabus event
       ev, err := uvalibrabus.MakeBusEvent(mbEvent.Detail)
	   if err != nil {
          fmt.Printf("ERROR: unmarshaling libra bus event (%s)", err.Error())
		  return err
	   }

       fmt.Printf("EVENT %s from:%s -> %s\n", mbEvent.ID, mbEvent.Source, ev.String())
    }
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}

//
// end of file
//