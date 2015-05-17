package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/darnould/linux-monitoring-cloudwatch/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/aws"
	"github.com/darnould/linux-monitoring-cloudwatch/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/service/cloudwatch"
)

func main() {
	ns := flag.String("namespace", "", "CloudWatch metric namespace (required)")
	reg := flag.String("region", "", "AWS Region")

	flag.Parse()

	if *ns == "" {
		fmt.Fprintln(os.Stderr, "Usage: ", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *reg == "" {
		r, err := region()
		if err != nil {
			log.Fatal("Can't find region: ", err)
		}
		reg = &r
	}

	mem, err := memoryUsage()
	if err != nil {
		log.Fatal("Can't get memory usage: ", err)
	}

	err = putMetric("MemoryUsage", "Percent", mem, *ns, *reg)
	if err != nil {
		log.Fatal("Can't put memory usage metric: ", err)
	}
}

func putMetric(name, unit string, value float64, namespace, region string) error {
	svc := cloudwatch.New(&aws.Config{Region: region})

	metric_input := &cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				MetricName: aws.String(name),
				Unit:       aws.String(unit),
				Value:      aws.Double(value),
			},
		},
		Namespace: aws.String(namespace),
	}

	_, err := svc.PutMetricData(metric_input)
	if awserr := aws.Error(err); awserr != nil {
		return fmt.Errorf("[%s] %s", awserr.Code, awserr.Message)
	} else if err != nil {
		return err
	}

	return nil
}

func memoryUsage() (percentUsed float64, err error) {
	b, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, fmt.Errorf("can't read memory stats: %s", err)
	}

	l := strings.Split(string(b), "\n")[0:2]
	totalLine, freeLine := l[0], l[1]

	re := regexp.MustCompile(`\s+(\d+)\s+`)

	t := re.FindStringSubmatch(totalLine)
	if len(t[1:]) < 1 {
		return 0, fmt.Errorf("can't parse total memory: %s", err)
	}
	total, err := strconv.ParseFloat(t[1], 32)
	if err != nil {
		return 0, fmt.Errorf("can't parse total memory: %s", err)
	}

	f := re.FindStringSubmatch(freeLine)
	if len(f[1:]) < 1 {
		return 0, fmt.Errorf("can't parse free memory: %s", err)
	}
	free, err := strconv.ParseFloat(f[1], 32)
	if err != nil {
		return 0, fmt.Errorf("can't parse free memory: %s", err)
	}

	used := total - free
	percentUsed = used / total * 100

	log.Print("Memory usage: ", percentUsed)

	return percentUsed, err
}

func region() (region string, err error) {
	resp, err := http.Get("http://169.254.169.254/latest/dynamic/instance-identity/document")
	if err != nil {
		return "", fmt.Errorf("can't reach metadata endpoint - %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("can't read metadata response body - %s", err)
	}

	var data map[string]string
	json.Unmarshal(body, &data)

	return data["region"], err
}
