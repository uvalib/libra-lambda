//
// main for lambda deployable
//

// include this on a lambda build only
//go:build lambda

package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return process(request.RequestContext.RequestID, "", request)
}

func main() {
	lambda.Start(HandleRequest)
}

//
// end of file
//
