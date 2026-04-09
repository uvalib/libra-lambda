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
	"path"
	"path/filepath"
	"strings"

	"github.com/uvalib/easystore/uvaeasystore"
	librametadata "github.com/uvalib/libra-metadata"
)

var metadataFilename = "metadata.json"
var fieldsFilename = "fields.json"
var auditFilename = "audit.json"
var descriptionFileName = "aptrust-description.txt"
var titleFileName = "aptrust-title.txt"
var manifestFilename = "manifest-md5.txt"

func createBagContent(cfg *Config, httpClient *http.Client, obj uvaeasystore.EasyStoreObject) (string, []string, error) {

	// create the bag name and working directory
	bagName := strings.Replace(cfg.BagNameTemplate, "{:oid}", obj.Id(), 1)
	workDir := filepath.Join(cfg.ScratchFilesystem, bagName)
	files := make([]string, 0)

	// clean the scratch filesystem (it can persist across lambda executions, who knew)
	err := cleanScratchFilesystem(cfg.ScratchFilesystem)
	if err != nil {
		fmt.Printf("ERROR: cleaning scratch filesystem [%s] (%s)\n", cfg.ScratchFilesystem, err.Error())
		return bagName, files, err
	}

	// create the working directory
	err = os.MkdirAll(workDir, 0755)
	if err != nil {
		fmt.Printf("ERROR: creating work directory [%s] (%s)\n", workDir, err.Error())
		return bagName, files, err
	}

	var buf []byte

	// if the metadata exists, write it
	if obj.Metadata() != nil {
		buf, err = obj.Metadata().Payload()
		if err != nil {
			fmt.Printf("ERROR: getting metadata payload (%s)\n", err.Error())
			return bagName, files, err
		}

		err = writeFile(filepath.Join(workDir, metadataFilename), buf)
		if err != nil {
			return bagName, files, err
		}

		// add to the bag files list
		files = append(files, metadataFilename)

		// write the title and description files
		meta, err := librametadata.ETDWorkFromBytes(buf)
		if err != nil {
			fmt.Printf("ERROR: creating libra metadata (%s)\n", err.Error())
			return bagName, files, err
		}

		if len(meta.Title) != 0 {
			err = writeFile(filepath.Join(workDir, titleFileName), []byte(meta.Title))
			if err != nil {
				return bagName, files, err
			}
			// add to the bag files list
			files = append(files, titleFileName)
		}

		if len(meta.Abstract) != 0 {
			err = writeFile(filepath.Join(workDir, descriptionFileName), []byte(meta.Abstract))
			if err != nil {
				return bagName, files, err
			}
			// add to the bag files list
			files = append(files, descriptionFileName)
		}
	}

	// if the fields exist, write them
	if obj.Fields() != nil {
		buf, err = json.Marshal(obj.Fields())
		if err != nil {
			fmt.Printf("ERROR: getting fields payload (%s)\n", err.Error())
			return bagName, files, err
		}

		err = writeFile(filepath.Join(workDir, fieldsFilename), buf)
		if err != nil {
			return bagName, files, err
		}

		// add to the bag files list
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
				return bagName, files, err
			}

			err = writeFile(filepath.Join(workDir, f.Name()), buf)
			if err != nil {
				return bagName, files, err
			}

			// add to the bag files list
			files = append(files, f.Name())
		}
	}

	// generate the audit
	files, err = generateAudit(cfg, httpClient, obj, workDir, files)
	if err != nil {
		return bagName, files, err
	}

	// generate the manifest
	files, err = generateManifest(workDir, files)

	// and we are done...
	return bagName, files, err
}

func generateManifest(workDir string, files []string) ([]string, error) {

	md5Data := ""
	for _, f := range files {
		fname := filepath.Join(workDir, f)
		fp, err := md5Checksum(fname)
		if err != nil {
			return files, err
		}
		md5Data += fmt.Sprintf("%s %s\n", fp, f)
	}

	// add to the bag files list
	files = append(files, manifestFilename)

	fname := filepath.Join(workDir, manifestFilename)
	return files, os.WriteFile(fname, []byte(md5Data), 0644)
}

func generateAudit(cfg *Config, httpClient *http.Client, obj uvaeasystore.EasyStoreObject, workDir string, files []string) ([]string, error) {

	// generate the query URL
	url := cfg.AuditQuery
	url = strings.Replace(url, "{:ns}", obj.Namespace(), 1)
	url = strings.Replace(url, "{:oid}", obj.Id(), 1)

	buf, err := httpGet(httpClient, url)
	// lets ignore errors for now
	if err != nil {
		fmt.Printf("WARNING: getting work audit information (%s)\n", err.Error())
		//		return files, err
		return files, nil
	}

	err = writeFile(filepath.Join(workDir, auditFilename), buf)
	if err != nil {
		return files, err
	}

	// add to the bag files list
	files = append(files, auditFilename)
	return files, err
}

func md5Checksum(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("ERROR: reading [%s] (%s)\n", filename, err.Error())
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum(data)), nil
}

func writeFile(filename string, buffer []byte) error {

	err := os.WriteFile(filename, buffer, 0644)
	if err != nil {
		fmt.Printf("ERROR: writing [%s] (%s)\n", filename, err.Error())
		return err
	}
	return nil
}

func cleanScratchFilesystem(dir string) error {

	de, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, d := range de {
		err = os.RemoveAll(path.Join(dir, d.Name()))
		if err != nil {
			return err
		}
	}
	
	return nil
}

//
// end of file
//
