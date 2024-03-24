package healthevent

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type EventTypeCategory string

const (
	EventTypeCategoryIssue               EventTypeCategory = "issue"
	EventTypeCategoryAccountNotification EventTypeCategory = "accountNotification"
	EventTypeCategoryInvestigation       EventTypeCategory = "investigation"
	EventTypeCategoryScheduledChange     EventTypeCategory = "scheduledChange"
)

type EventScopeCode string

const (
	EventScopeCodeAccountSpecific EventScopeCode = "ACCOUNT_SPECIFIC"
	EventScopeCodePublic          EventScopeCode = "PUBLIC"
)

type StatusCode string

const (
	StatusCodeIssueOpen                 StatusCode = "open"
	StatusCodeIssueClosed               StatusCode = "closed"
	StatusCodeIssueUpcoming             StatusCode = "upcoming"
	StatusCodeScheduledChangesUpcoming  StatusCode = "Upcoming"
	StatusCodeScheduledChangesOngoing   StatusCode = "Ongoing"
	StatusCodeScheduledChangesCompleted StatusCode = "Completed"
	StatusCodeUndefined                 StatusCode = "-"
)

// https://docs.aws.amazon.com/health/latest/ug/cloudwatch-events-health.html#aws-health-event-schema
type Detail struct {
	EventArn          string                 `json:"eventArn"`
	Service           string                 `json:"service"`
	EventTypeCode     string                 `json:"eventTypeCode"`
	EventTypeCategory EventTypeCategory      `json:"eventTypeCategory"`
	EventScopeCode    EventScopeCode         `json:"eventScopeCode"`
	CommunicationID   string                 `json:"communicationId"`
	StartTime         *Time                  `json:"startTime"`
	EndTime           *Time                  `json:"endTime"`
	LastUpdatedTime   *Time                  `json:"lastUpdatedTime"`
	StatusCode        StatusCode             `json:"statusCode"`
	EventRegion       string                 `json:"eventRegion"`
	EventDescription  []*EventDescriptionRow `json:"eventDescription"`
	EventMetadata     map[string]string      `json:"eventMetadata"`
	AffectedEntities  []*AffectedEntity      `json:"affectedEntities"`
	Page              string                 `json:"page"`
	TotalPages        string                 `json:"totalPages"`
	AffectedAccount   string                 `json:"affectedAccount"`
}

type EventDescriptionRow struct {
	Language          string `json:"language"`
	LatestDescription string `json:"latestDescription"`
}

type AffectedEntity struct {
	EntityValue     string `json:"entityValue"`
	LastUpdatedTime *Time  `json:"lastUpdatedtime"`
	Status          string `json:"status"`
}

func ParseDetail(e events.CloudWatchEvent) (*Detail, error) {
	eventDetail := &Detail{}
	if err := json.Unmarshal(e.Detail, eventDetail); err != nil {
		return nil, err
	}
	return eventDetail, nil
}

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		t = nil
		return
	}
	// The format is similar to RFC1123 but different.
	// DoW, DD MMM YYYY HH:MM:SS TZ
	t.Time, err = time.Parse("Mon, _2 Jan 2006 15:04:05 MST", s)
	return
}
