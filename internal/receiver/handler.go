package receiver

import (
	"context"
	"log"

	"github.com/Arthur1/mackerel-aws-health-events-notifier/healthevent"
	"github.com/aws/aws-lambda-go/events"
)

type Handler struct{}

func New() *Handler {
	return new(Handler)
}

func (h *Handler) Handle(_ context.Context, e events.CloudWatchEvent) {
	eventDetail, err := healthevent.ParseDetail(e)
	if err != nil {
		panic(err)
	}
	log.Printf("%#v\n", eventDetail)
}
