package loganalytics

import (
	"testing"
	"time"

	azhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/open-telemetry/opentelemetry-log-collection/entry"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	var (
		operator             LogAnalyticsInput
		testPartitionKey     string    = "1"
		testSequenceNum      int64     = 600
		testTime             time.Time = time.Now()
		testOffset           int64     = 2000
		testTimeGeneratedSTR string    = "2021-04-26T17:11:41.3500000Z"
	)

	testTimeGenerated, err := time.Parse("2006-01-02T15:04:05.0000000Z07", testTimeGeneratedSTR)
	require.NoError(t, err)

	cases := []struct {
		name     string
		event    azhub.Event
		records  map[string]interface{}
		expected *entry.Entry
	}{
		{
			"promote-type-and-timegenerated",
			azhub.Event{
				Data:         nil,
				ID:           "000",
				PartitionKey: &testPartitionKey,
				Properties: map[string]interface{}{
					"user": "stanza",
					"id":   1,
					"root": false,
				},
				SystemProperties: &azhub.SystemProperties{
					SequenceNumber: &testSequenceNum,
					EnqueuedTime:   &testTime,
					Offset:         &testOffset,
				},
			},
			map[string]interface{}{
				"aks_cluster":   "us-east1-b-dev-0",
				"system_id":     100,
				"dev":           false,
				"timegenerated": testTimeGeneratedSTR,
				"type":          "unit_test",
			},
			&entry.Entry{
				Timestamp: testTimeGenerated,
				Attributes: map[string]string{
					"azure_log_analytics_table": "unit_test",
				},
				Body: map[string]interface{}{
					"partition_key": &testPartitionKey,
					"properties": map[string]interface{}{
						"user": "stanza",
						"id":   1,
						"root": false,
					},
					"system_properties": map[string]interface{}{
						"x-opt-sequence-number": &testSequenceNum,
						"x-opt-enqueued-time":   &testTime,
						"x-opt-offset":          &testOffset,
					},
					"unit_test": map[string]interface{}{
						"aks_cluster":   "us-east1-b-dev-0",
						"system_id":     100,
						"dev":           false,
						"timegenerated": testTimeGeneratedSTR,
					},
				},
				Resource: map[string]string{
					"event_id": "000",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := entry.New()
			err := operator.parse(tc.event, tc.records, e)
			require.NoError(t, err)
			require.Equal(t, tc.expected, e)
		})
	}
}
