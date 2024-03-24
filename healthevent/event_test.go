package healthevent

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestParseDetail(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		eventFile string
		want      *Detail
	}{
		// https://docs.aws.amazon.com/health/latest/ug/cloudwatch-events-health.html#amazon-ec2-operational-issue
		"Public Health Event - Amazon EC2 operational issue": {
			eventFile: "testdata/public-ec2-event.json",
			want: &Detail{
				EventArn:          "arn:aws:health:af-south-1::event/EC2/AWS_EC2_OPERATIONAL_ISSUE/AWS_EC2_OPERATIONAL_ISSUE_7f35c8ae-af1f-54e6-a526-d0179ed6d68f",
				Service:           "EC2",
				EventTypeCode:     "AWS_EC2_OPERATIONAL_ISSUE",
				EventTypeCategory: EventTypeCategoryIssue,
				EventScopeCode:    EventScopeCodePublic,
				CommunicationID:   "01b0993207d81a09dcd552ebd1e633e36cf1f09a-1",
				StartTime:         strToTime(t, "2023-01-27T06:02:51Z"),
				EndTime:           strToTime(t, "2023-01-27T09:01:22Z"),
				LastUpdatedTime:   strToTime(t, "2023-01-27T09:01:22Z"),
				StatusCode:        StatusCodeIssueOpen,
				EventRegion:       "af-south-1",
				EventDescription: []*EventDescriptionRow{{
					Language:          "en_US",
					LatestDescription: "Current severity level: Operating normally\n\n[RESOLVED] \n\n [03:15 PM PST] We continue see recovery \n\nThe following AWS services were previously impacted but are now operating normally: APPSYNC, BACKUP, EVENTS.",
				}},
				AffectedEntities: []*AffectedEntity{},
				Page:             "1",
				TotalPages:       "1",
				AffectedAccount:  "123456789012",
			},
		},

		// https://docs.aws.amazon.com/health/latest/ug/cloudwatch-events-health.html#elastic-load-balancing-api-issue
		"Account-specific AWS Health Event - Elastic Load Balancing API Issue": {
			eventFile: "testdata/specific-elb-event.json",
			want: &Detail{
				EventArn:          "arn:aws:health:ap-southeast-2::event/AWS_ELASTICLOADBALANCING_API_ISSUE_90353408594353980",
				Service:           "ELASTICLOADBALANCING",
				EventTypeCode:     "AWS_ELASTICLOADBALANCING_API_ISSUE",
				EventTypeCategory: EventTypeCategoryIssue,
				EventScopeCode:    EventScopeCodeAccountSpecific,
				CommunicationID:   "01b0993207d81a09dcd552ebd1e633e36cf1f09a-1",
				StartTime:         strToTime(t, "2022-06-10T05:01:10Z"),
				EndTime:           strToTime(t, "2022-06-10T05:30:57Z"),
				StatusCode:        StatusCodeIssueOpen,
				EventRegion:       "ap-southeast-2",
				EventDescription: []*EventDescriptionRow{{
					Language:          "en_US",
					LatestDescription: "A description of the event will be provided here",
				}},
				Page:            "1",
				TotalPages:      "1",
				AffectedAccount: "123456789012",
			},
		},

		// https://docs.aws.amazon.com/health/latest/ug/cloudwatch-events-health.html#amazon-ec2-instance-store-drive-performance-degraded
		"Account-specific AWS Health Event - Amazon EC2 Instance Store Drive Performance Degraded": {
			eventFile: "testdata/specific-ec2-event.json",
			want: &Detail{
				EventArn:          "arn:aws:health:us-west-2::event/AWS_EC2_INSTANCE_STORE_DRIVE_PERFORMANCE_DEGRADED_90353408594353980",
				Service:           "EC2",
				EventTypeCode:     "AWS_EC2_INSTANCE_STORE_DRIVE_PERFORMANCE_DEGRADED",
				EventTypeCategory: EventTypeCategoryIssue,
				EventScopeCode:    EventScopeCodeAccountSpecific,
				CommunicationID:   "01b0993207d81a09dcd552ebd1e633e36cf1f09a-1",
				StartTime:         strToTime(t, "2022-06-03T05:01:10Z"),
				EndTime:           strToTime(t, "2022-06-03T05:30:57Z"),
				StatusCode:        StatusCodeIssueOpen,
				EventRegion:       "us-west-2",
				EventDescription: []*EventDescriptionRow{{
					Language:          "en_US",
					LatestDescription: "A description of the event will be provided here",
				}},
				AffectedEntities: []*AffectedEntity{{EntityValue: "i-abcd1111"}},
				Page:             "1",
				TotalPages:       "1",
				AffectedAccount:  "123456789012",
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := events.CloudWatchEvent{}
			bytes, err := os.ReadFile(tt.eventFile)
			assert.NoError(t, err)
			err = json.Unmarshal(bytes, &e)
			assert.NoError(t, err)

			got, err := ParseDetail(e)
			assert.NoError(t, err)
			assert.Empty(t, cmp.Diff(tt.want, got))
		})
	}
}

func strToTime(t *testing.T, str string) *Time {
	ti, err := time.Parse(time.RFC3339, str)
	assert.NoError(t, err)
	return &Time{ti}
}
