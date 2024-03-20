//
//
//

// include this on a cmdline build only
//go:build cmdline

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	var messageId string
	var source string

	flag.StringVar(&messageId, "messageid", "0-0-0-0", "Message identifier")
	flag.StringVar(&source, "source", "the.source", "Message source")
	flag.Parse()

	pl := []byte("{}")
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
