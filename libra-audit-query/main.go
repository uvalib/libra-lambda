//
//
//

// include this on a cmdline build only
//go:build cmdline

package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"os"
)

func main() {

	var messageId string
	var namespace string
	var objectId string
	var who string

	flag.StringVar(&messageId, "messageid", "0-0-0-0", "Message identifier")
	flag.StringVar(&namespace, "namespace", "", "Object namespace")
	flag.StringVar(&objectId, "objid", "", "Object identifier")
	flag.StringVar(&who, "who", "", "Who (computing id)")
	flag.Parse()

	if (len(namespace) == 0 || len(objectId) == 0) && len(who) == 0 {
		fmt.Printf("ERROR: incorrect commandline, use --help for details\n")
		os.Exit(1)
	}

	req := events.APIGatewayProxyRequest{}
	req.QueryStringParameters = map[string]string{}
	if len(who) != 0 {
		req.QueryStringParameters["who"] = who
	} else {
		req.QueryStringParameters["namespace"] = namespace
		req.QueryStringParameters["oid"] = objectId
	}

	resp, err := process(messageId, "api.gateway", req)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("INFO: response: %s\n", resp.Body)
	fmt.Printf("INFO: terminating with HTTP %d\n", resp.StatusCode)
}

//
// end of file
//
