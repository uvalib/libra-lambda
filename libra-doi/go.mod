module github.com/uvalib/libra-doi

go 1.25.0

require (
	github.com/aws/aws-lambda-go v1.54.0
	github.com/davecgh/go-spew v1.1.1
	github.com/uvalib/easystore/uvaeasystore v0.0.0-20260413184000-ac1e96bfa2b7
	github.com/uvalib/libra-metadata v0.0.0-20250513131340-aa4ee04ad7d1
	github.com/uvalib/librabus-sdk/uvalibrabus v0.0.0-20260406142030-486f51674d88
)

// for local development
//require github.com/uvalib/easystore/uvaeasystore v0.0.0
//replace github.com/uvalib/easystore/uvaeasystore => ../../easystore/uvaeasystore

require (
	github.com/aws/aws-sdk-go-v2 v1.41.6 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.9 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.32.16 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.15 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.22 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.22.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatchevents v1.32.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.100.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.20 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.42.0 // indirect
	github.com/aws/smithy-go v1.25.1 // indirect
	github.com/lib/pq v1.12.3 // indirect
	github.com/rs/xid v1.6.0 // indirect
	golang.org/x/exp v0.0.0-20260410095643-746e56fc9e2f // indirect
)
