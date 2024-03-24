package main

import (
	"github.com/Arthur1/mackerel-aws-health-events-notifier/internal/receiver"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	h := receiver.New()
	lambda.Start(h.Handle)
}
