//
//
//

package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/uvalib/easystore/uvaeasystore"
)

var metadataFilename = "metadata.json"
var fieldsFilename = "fields.json"
var manifestFilename = "manifest-md5.txt"
var descriptionFileName = "aptrust-description.txt"
var titleFileName = "aptrust-title.txt"

func createBagContent(cfg *Config, httpClient *http.Client, obj uvaeasystore.EasyStoreObject) (string, error) {

	// create the bag name and working directory
	bagName := strings.Replace(cfg.BagNameTemplate, "{:oid}", obj.Id(), 1)
	workDir := filepath.Join(cfg.ScratchFilesystem, bagName)

	// create the working directory
	err := os.MkdirAll(workDir, 0755)
	if err != nil {
		fmt.Printf("ERROR: creating work directory [%s] (%s)\n", workDir, err.Error())
		return bagName, err
	}

	var buf []byte
	files := make([]string, 0)

	// if the metadata exists, write it
	if obj.Metadata() != nil {
		buf, err = obj.Metadata().Payload()
		if err != nil {
			fmt.Printf("ERROR: getting metadata payload (%s)\n", err.Error())
			return bagName, err
		}

		fname := filepath.Join(workDir, metadataFilename)
		err = os.WriteFile(fname, buf, 0644)
		if err != nil {
			fmt.Printf("ERROR: writing [%s] (%s)\n", fname, err.Error())
			return bagName, err
		}

		files = append(files, metadataFilename)

		// write the title and description files

	}

	// if the fields exist, write them
	if obj.Fields() != nil {
		buf, err = json.Marshal(obj.Fields())
		if err != nil {
			fmt.Printf("ERROR: getting fields payload (%s)\n", err.Error())
			return bagName, err
		}

		fname := filepath.Join(workDir, fieldsFilename)
		err = os.WriteFile(fname, buf, 0644)
		if err != nil {
			fmt.Printf("ERROR: writing [%s] (%s)\n", fname, err.Error())
			return bagName, err
		}
		files = append(files, fieldsFilename)
	}

	// if the files exist, write them
	if obj.Files() != nil {
		for _, f := range obj.Files() {
			if len(f.Url()) != 0 {
				buf, err = httpGet(httpClient, f.Url())
			} else {
				buf, err = f.Payload()
			}
			if err != nil {
				fmt.Printf("ERROR: getting file payload (%s)\n", err.Error())
				return bagName, err
			}

			fname := filepath.Join(workDir, f.Name())
			err = os.WriteFile(fname, buf, 0644)
			if err != nil {
				fmt.Printf("ERROR: writing [%s] (%s)\n", fname, err.Error())
				return bagName, err
			}
			files = append(files, f.Name())
		}
	}

	// generate the manifest
	err = generateManifest(workDir, files)

	// and done...
	return bagName, err
}

func generateManifest(workDir string, files []string) error {

	md5Data := ""
	for _, f := range files {
		fname := filepath.Join(workDir, f)
		fp, err := md5Checksum(fname)
		if err != nil {
			return err
		}
		md5Data += fmt.Sprintf("%s %s\n", fp, f)
	}

	fname := filepath.Join(workDir, manifestFilename)
	return os.WriteFile(fname, []byte(md5Data), 0644)
}

func md5Checksum(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("ERROR: reading [%s] (%s)\n", filename, err.Error())
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum(data)), nil
}

//
// end of file
//
