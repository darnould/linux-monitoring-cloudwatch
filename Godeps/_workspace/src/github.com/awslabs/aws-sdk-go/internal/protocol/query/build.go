package query

//go:generate go run ../../fixtures/protocol/generate.go ../../fixtures/protocol/input/query.json build_test.go

import (
	"net/url"

	"github.com/darnould/linux-monitoring-cloudwatch/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/aws"
	"github.com/darnould/linux-monitoring-cloudwatch/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/internal/protocol/query/queryutil"
)

func Build(r *aws. // Build builds a request for an AWS Query service.
Request) {
	body := url.Values{
		"Action":  {r.Operation.Name},
		"Version": {r.Service.APIVersion},
	}
	if err := queryutil.Parse(body, r.Params, false); err != nil {
		r.Error = err
		return
	}

	if r.ExpireTime == 0 {
		r.HTTPRequest.Method = "POST"
		r.HTTPRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		r.SetBufferBody([]byte(body.Encode()))
	} else { // This is a pre-signed request
		r.HTTPRequest.Method = "GET"
		r.HTTPRequest.URL.RawQuery = body.Encode()
	}
}
