//
//
//

// include this on a lambda build only
//go:build lambda

package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event events.EventBridgeEvent) error {

	// process the message, in the event of an error, it is re-queued
	return process(event.ID, event.Source, event.Detail)
}

func main() {
	lambda.Start(HandleRequest)
}

//
// end of file
//
