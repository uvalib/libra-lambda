//
// simple module to get and set parameter values in the ssm
//

package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client

func initS3() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	s3Client = s3.NewFromConfig(cfg)
	return nil
}

func putS3(bucket string, key string, buffer []byte) error {

	fmt.Printf("DEBUG: uploading s3://%s/%s\n", bucket, key)

	_, err := s3Client.PutObject(context.TODO(),
		&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   bytes.NewReader(buffer),
		})

	if err != nil {
		return err
	}

	return nil
}

//
// end of file
//
