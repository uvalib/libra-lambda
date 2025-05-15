//
//
//

package main

import (
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	"net/http"
	"strings"
)

func notifySis(config *Config, fields uvaeasystore.EasyStoreObjectFields, auth string, client *http.Client) error {

	sisId := strings.Replace(fields["source-id"], "sis:", "", 1)
	doi := fields["doi"]

	// substitute values into url
	url := strings.Replace(config.SisNotifyUrl, "{:id}", sisId, 1)
	url = strings.Replace(url, "{:auth}", auth, 1)
	url = strings.Replace(url, "{:doi}", doi, 1)

	buf, err := httpPut(client, url, nil, "")
	if err != nil {
		if buf != nil {
			fmt.Printf("ERROR: failed response [%s]\n", string(buf))
		}
		return err
	}

	return nil
}

//
// end of file
//
