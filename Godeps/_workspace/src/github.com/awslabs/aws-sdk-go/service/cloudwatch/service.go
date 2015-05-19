package cloudwatch

import (
	"github.com/darnould/linux-monitoring-cloudwatch/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/aws"
	"github.com/darnould/linux-monitoring-cloudwatch/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/internal/protocol/query"
	"github.com/darnould/linux-monitoring-cloudwatch/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/internal/signer/v4"
)

type CloudWatch struct { // CloudWatch is a client for CloudWatch.
	*aws.Service
}

// Used for custom service initialization logic
var initService func(*aws.Service)

// Used for custom request initialization logic
var initRequest func(*aws.Request)

// New returns a new CloudWatch client.
func New(config *aws.Config) *CloudWatch {
	if config == nil {
		config = &aws.Config{}
	}

	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "monitoring",
		APIVersion:  "2010-08-01",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(query.Build)
	service.Handlers.Unmarshal.PushBack(query.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(query.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(query.UnmarshalError)

	// Run custom service initialization if present
	if initService != nil {
		initService(service)
	}

	return &CloudWatch{service}
}

// newRequest creates a new request for a CloudWatch operation and runs any
// custom request initialization.
func (c *CloudWatch) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
