<h1 align="center">linux-monitoring-cloudwatch</h1>

<p align="center">
  Send <code>MemoryUsage</code> percentage metrics to AWS CloudWatch from Linux-based EC2 instances.
</p>

<p align="center">
  <a href="https://travis-ci.org/darnould/linux-monitoring-cloudwatch" target="_blank"><img src="https://img.shields.io/travis/darnould/linux-monitoring-cloudwatch.svg?style=flat-square"></a>
</p>

### Why is this useful?

EC2 natively only sends hypervisor-visible metrics to CloudWatch (e.g. CPU usage).  Others such as memory usage require OS-level integration - this is where `linux-monitoring-cloudwatch` comes in.

No further dependencies are necessary on your instance - just the executable.

With the metrics emitted by `linux-monitoring-cloudwatch` (`MemoryUsage` at this time), you may trigger alarms & autoscaling.

### Example usage
```sh
./linux-monitoring-cloudwatch --namespace my-app/int --region eu-west-1
```

`--region` is optional: without providing it, `linux-monitoring-cloudwatch` will query the instance metadata API at runtime.

### Building (with Docker)
```sh
docker run -v "$PWD":/go/src/github.com/darnould/linux-monitoring-cloudwatch -w /go/src/github.com/darnould/linux-monitoring-cloudwatch golang:1.3 go build
```
This will produce a single `linux-monitoring-cloudwatch` executable.

### Installation
 * [Build](#building-with-docker) or [download](https://github.com/darnould/linux-monitoring-cloudwatch/releases/download/v0.1.0/linux-monitoring-cloudwatch) the executable.
 * Copy the `linux-monitoring-cloudwatch` executable to your instances' filesystems.
 * Allow your instances to make `PutMetricData` calls to CloudWatch with an IAM policy such as:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "cloudwatch:PutMetricData",
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
```
* A cron job such as the following will report metrics to CloudWatch every 5 minutes:
```
/5 * * * *  /path/to/linux-monitoring-cloudwatch --namespace NAMESPACE --region REGION &>> /var/log/linux-monitoring-cloudwatch.log
```

### Coming up
  * `DiskUsage` percentage metric
