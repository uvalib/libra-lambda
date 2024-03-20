//
// simple module to get and set parameter values in the ssm
//

package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var ssmClient *ssm.Client

func initParameter() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	ssmClient = ssm.NewFromConfig(cfg)
	return nil
}

func getParameter(name string) (string, error) {
	param, err := ssmClient.GetParameter(context.Background(),
		&ssm.GetParameterInput{
			Name:           aws.String(name),
			WithDecryption: aws.Bool(false),
		})

	if err != nil {
		return "", err
	}

	return *param.Parameter.Value, nil
}

func setParameter(name string, value string) error {
	_, err := ssmClient.PutParameter(context.Background(),
		&ssm.PutParameterInput{
			Name:      aws.String(name),
			Value:     aws.String(value),
			Overwrite: aws.Bool(true),
		})

	if err != nil {
		return err
	}

	return nil
}

//
// end of file
//
