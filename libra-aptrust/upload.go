//
//
//

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func uploadContent(cfg *Config, s3 *s3.Client, bucket string, prefix string, bagName string, files []string) error {

	// this is our content directory
	contentDir := filepath.Join(cfg.ScratchFilesystem, bagName)

	// create a new uploader
	uploader := manager.NewUploader(s3)

	fullPrefix := filepath.Join(prefix, bagName)
	for _, fn := range files {
		remoteName := filepath.Join(fullPrefix, fn)
		localName := filepath.Join(contentDir, fn)
		err := uploadFile(uploader, bucket, remoteName, localName)
		if err != nil {
			return err
		}
	}

	return nil
}

func uploadFile(uploader *manager.Uploader, bucket string, key string, local string) error {

	target := fmt.Sprintf("s3://%s/%s", bucket, key)
	fmt.Printf("DEBUG: put from %s to %s\n", local, target)

	// open the file
	file, err := os.Open(local)
	if err != nil {
		// assume the error is file not found... probably reasonable
		return os.ErrNotExist
	}
	defer file.Close()

	// get the filesize
	s, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := s.Size()

	// Upload the file to S3.
	start := time.Now()
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	})
	if err != nil {
		return err
	}

	duration := time.Since(start)
	fmt.Printf("DEBUG: put %s complete in %0.2f seconds (%d bytes, %0.2f bytes/sec)\n", target, duration.Seconds(), fileSize, float64(fileSize)/duration.Seconds())
	return nil
}

func fileList(dirname string) ([]string, error) {

	files, err := os.ReadDir(dirname)
	if err != nil {
		fmt.Printf("ERROR: reading %s (%s)\n", dirname, err.Error())
		return nil, err
	}

	results := make([]string, 0)
	for _, file := range files {
		if file.IsDir() == false {
			results = append(results, file.Name())
		}
	}
	return results, nil
}

//
// end of file
//
