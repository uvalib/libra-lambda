module github.com/uvalib/libra-doi

go 1.23.0

toolchain go1.24.2

require (
	github.com/aws/aws-lambda-go v1.49.0
	github.com/davecgh/go-spew v1.1.1
	github.com/uvalib/easystore/uvaeasystore v0.0.0-20250717150932-d438f30cb70b
	github.com/uvalib/libra-metadata v0.0.0-20250513131340-aa4ee04ad7d1
	github.com/uvalib/librabus-sdk/uvalibrabus v0.0.0-20250722185634-5799610d7eb7
)

// for local development
//require github.com/uvalib/easystore/uvaeasystore v0.0.0
//replace github.com/uvalib/easystore/uvaeasystore => ../../easystore/uvaeasystore

require (
	github.com/aws/aws-sdk-go-v2 v1.36.6 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.11 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.29.18 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.71 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.33 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.17.85 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.37 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.37 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.37 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatchevents v1.28.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.84.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.34.1 // indirect
	github.com/aws/smithy-go v1.22.4 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/rs/xid v1.6.0 // indirect
	golang.org/x/exp v0.0.0-20250718183923-645b1fa84792 // indirect
)
