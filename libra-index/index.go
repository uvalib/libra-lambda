//
//
//

package main

import (
	"encoding/json"
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	"net/http"
)

type IndexWork struct {
	Id       string          `json:"id,omitempty"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
	Fields   json.RawMessage `json:"fields,omitempty"`
	Files    json.RawMessage `json:"files,omitempty"`
}

func updateIndex(config *Config, eso uvaeasystore.EasyStoreObject, client *http.Client) error {

	// create the request payload
	req := IndexWork{}
	req.Id = eso.Id()

	// include metadata if it exists
	if eso.Metadata() != nil {
		buf, err := eso.Metadata().Payload()
		if err != nil {
			return err
		}
		req.Metadata = buf
	}

	// include fields if they exist
	if eso.Fields() != nil {
		buf, err := json.Marshal(eso.Fields())
		if err != nil {
			return err
		}
		req.Fields = buf
	}

	// include file details if they exist
	if eso.Files() != nil {
		files := make([]string, 0, 10)
		for _, f := range eso.Files() {
			files = append(files, f.Name())
		}
		buf, err := json.Marshal(files)
		if err != nil {
			return err
		}
		req.Files = buf
	}

	pl, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("ERROR: json marshal of IndexWork (%s)\n", err.Error())
		return err
	}

	fmt.Printf("INFO: payload [%s]\n", string(pl))

	buf, err := httpPost(client, config.IndexUpdateUrl, pl, "application/json")
	if err != nil {
		fmt.Printf("ERROR: failed payload [%s]\n", string(pl))
		if buf != nil {
			fmt.Printf("ERROR: failed response [%s]\n", string(buf))
		}
		return err
	}

	fmt.Printf("INFO: response [%s]\n", string(buf))

	// all good
	return nil
}

//
// end of file
//
