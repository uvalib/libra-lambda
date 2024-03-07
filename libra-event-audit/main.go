//
//
//

package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
)

func HandleRequest(ctx context.Context, event *events.EventBridgeEvent) (*string, error) {

	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	ev, err := uvalibrabus.MakeBusEvent(event.Detail)
	if err != nil {
		return nil, err
	}

	fmt.Printf("EVENT %s from:%s -> %s\n", event.ID, event.Source, ev.String())
	return nil, nil
}

func main() {
	lambda.Start(HandleRequest)
}

//
// end of file
//