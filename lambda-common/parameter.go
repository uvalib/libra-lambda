//
// simple module to get and set parameter values in the ssm
//

package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func newParameterClient() (*ssm.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	return ssm.NewFromConfig(cfg), nil
}

func getParameter(client *ssm.Client, name string) (string, error) {
	param, err := client.GetParameter(context.Background(),
		&ssm.GetParameterInput{
			Name:           aws.String(name),
			WithDecryption: aws.Bool(false),
		})

	if err != nil {
		fmt.Printf("ERROR: getting parameter (%s)\n", err.Error())
		return "", err
	}

	return *param.Parameter.Value, nil
}

func setParameter(client *ssm.Client, name string, value string) error {
	_, err := client.PutParameter(context.Background(),
		&ssm.PutParameterInput{
			Name:      aws.String(name),
			Value:     aws.String(value),
			Overwrite: aws.Bool(true),
		})

	if err != nil {
		fmt.Printf("ERROR: setting parameter (%s)\n", err.Error())
		return err
	}

	return nil
}

//
// end of file
//
