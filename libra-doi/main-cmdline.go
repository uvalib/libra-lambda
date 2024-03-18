//
//
//

// include this on a cmdline build only
//go:build cmdline

package main

import (
	"flag"
	"fmt"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
	"os"
)

func main() {

	var messageId string
	var source string
	var eventName string
	var namespace string
	var objectId string

	flag.StringVar(&messageId, "messageid", "0-0-0-0", "Message identifier")
	flag.StringVar(&source, "source", "the.source", "Message source")
	flag.StringVar(&eventName, "eventname", "", "Event name")
	flag.StringVar(&namespace, "namespace", "", "Object namespace")
	flag.StringVar(&objectId, "objid", "", "Object identifier")
	flag.Parse()

	if len(eventName) == 0 || len(namespace) == 0 || len(objectId) == 0 {
		fmt.Printf("ERROR: incorrect commandline, use --help for details\n")
		os.Exit(1)
	}

	ev := uvalibrabus.UvaBusEvent{}
	ev.EventName = eventName
	ev.Namespace = namespace
	ev.Identifier = objectId

	pl, _ := ev.Serialize()
	err := process(messageId, source, pl)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("INFO: terminating normally\n")
}

//
// end of file
//
