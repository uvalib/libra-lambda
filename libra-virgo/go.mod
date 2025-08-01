module github.com/uvalib/libra-virgo

go 1.24.3

require (
	github.com/aws/aws-lambda-go v1.49.0
	github.com/aws/aws-sdk-go-v2 v1.37.1
	github.com/aws/aws-sdk-go-v2/config v1.30.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.85.1
	github.com/uvalib/easystore/uvaeasystore v0.0.0-20250723164731-027ac39929ad
	github.com/uvalib/libra-metadata v0.0.0-20250513131340-aa4ee04ad7d1
	github.com/uvalib/librabus-sdk/uvalibrabus v0.0.0-20250801130056-157231a1fcac
)

// for local development
//require github.com/uvalib/easystore/uvaeasystore v0.0.0
//replace github.com/uvalib/easystore/uvaeasystore => ../../easystore/uvaeasystore

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.18.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.18.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatchevents v1.29.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.26.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.31.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.35.1 // indirect
	github.com/aws/smithy-go v1.22.5 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/rs/xid v1.6.0 // indirect
	golang.org/x/exp v0.0.0-20250718183923-645b1fa84792 // indirect
)
