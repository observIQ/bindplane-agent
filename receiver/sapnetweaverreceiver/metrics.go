// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sapnetweaverreceiver // import "github.com/observiq/observiq-otel-collector/receiver/sapnetweaverreceiver"

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/receiver/scrapererror"

	"github.com/observiq/observiq-otel-collector/receiver/sapnetweaverreceiver/internal/metadata"
	"github.com/observiq/observiq-otel-collector/receiver/sapnetweaverreceiver/internal/models"
)

const (
	collectMetricError              = "failed to collect metric %s: %w"
	collectMetricErrorWithAttribute = "failed to collect metric %s with attribute %s: %w"
	// MBToBytes converts 1 megabytes to byte
	MBToBytes = 1000000
)

var (
	errValueNotFound     = errors.New("value not found")
	errValueHyphen       = errors.New("'-' value found")
	errInvalidStateColor = errors.New("invalid control state color value")
)

func (s *sapNetweaverScraper) recordSapnetweaverWorkProcessesActiveCount(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "Total Number of Work Processes"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverWorkProcessesCountDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverSystemAvailabilityDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "Availability"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverSystemAvailabilityDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverSystemUtilizationDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "System Utilization"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverSystemUtilizationDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverMemoryUsageDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	// used this name for Percentage_Used which has a parent ID referencing Swap_Space
	metricName := "Swap_Space_Percentage_Used"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverMemoryUsageDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverMemoryConfiguredDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "Configured Memory"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	mbytes, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		errs.AddPartial(1, fmt.Errorf("failed to parse int64 for SapnetweaverMemoryConfigured, value was %v: %w", val, err))
		return
	}

	s.mb.RecordSapnetweaverMemoryConfiguredDataPoint(now, mbytes*int64(MBToBytes))
}

func (s *sapNetweaverScraper) recordSapnetweaverMemoryFreeDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "Free Memory"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	mbytes, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		errs.AddPartial(1, fmt.Errorf("failed to parse int64 for SapnetweaverMemoryFree, value was %v: %w", val, err))
		return
	}

	s.mb.RecordSapnetweaverMemoryFreeDataPoint(now, mbytes*int64(MBToBytes))
}

func (s *sapNetweaverScraper) recordSapnetweaverSessionCountDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "Number of Sessions"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverSessionCountDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverQueueCountDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "QueueLen"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverQueueCountDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverQueuePeakCountDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "PeakQueueLen"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverQueuePeakCountDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverJobAbortedDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "AbortedJobs"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverJobAbortedDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverAbapUpdateErrorsDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "AbapErrorInUpdate"
	val, ok := alertTreeResponse[metricName]
	if !ok {
		errs.AddPartial(1, fmt.Errorf(collectMetricError, metricName, errValueNotFound))
		return
	}

	switch models.StateColor(val) {
	case models.StateColorGray:
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 1, metadata.AttributeControlStateGrey)
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 0, metadata.AttributeControlStateGreen)
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 0, metadata.AttributeControlStateYellow)
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 0, metadata.AttributeControlStateRed)
	case models.StateColorGreen:
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 0, metadata.AttributeControlStateGrey)
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 1, metadata.AttributeControlStateGreen)
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 0, metadata.AttributeControlStateYellow)
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 0, metadata.AttributeControlStateRed)
	case models.StateColorYellow:
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 0, metadata.AttributeControlStateGrey)
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 0, metadata.AttributeControlStateGreen)
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 1, metadata.AttributeControlStateYellow)
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 0, metadata.AttributeControlStateRed)
	case models.StateColorRed:
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 0, metadata.AttributeControlStateGrey)
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 0, metadata.AttributeControlStateGreen)
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 0, metadata.AttributeControlStateYellow)
		s.mb.RecordSapnetweaverAbapUpdateErrorsDataPoint(now, 1, metadata.AttributeControlStateRed)
	default:
		errs.AddPartial(1, fmt.Errorf(collectMetricError, metricName, errInvalidStateColor))
	}
}

func (s *sapNetweaverScraper) RecordSapnetweaverResponseDurationDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	s.recordSapnetweaverResponseDurationDataPointDialog(now, alertTreeResponse, errs)
	s.recordSapnetweaverResponseDurationDataPointDialogRFC(now, alertTreeResponse, errs)
	s.recordSapnetweaverResponseDurationDataPointTransaction(now, alertTreeResponse, errs)
	s.recordSapnetweaverResponseDurationDataPointHTTP(now, alertTreeResponse, errs)
}

func (s *sapNetweaverScraper) recordSapnetweaverResponseDurationDataPointDialog(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "ResponseTimeDialog"
	val, err := parseResponse(metricName, metadata.AttributeResponseTypeDialog.String(), alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverResponseDurationDataPoint(now, val, metadata.AttributeResponseTypeDialog)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverResponseDurationDataPointDialogRFC(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "ResponseTimeDialogRFC"
	val, err := parseResponse(metricName, metadata.AttributeResponseTypeDialogRFC.String(), alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverResponseDurationDataPoint(now, val, metadata.AttributeResponseTypeDialogRFC)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverResponseDurationDataPointTransaction(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "ResponseTime(StandardTran.)"
	val, err := parseResponse(metricName, metadata.AttributeResponseTypeTransaction.String(), alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverResponseDurationDataPoint(now, val, metadata.AttributeResponseTypeTransaction)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverResponseDurationDataPointHTTP(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "ResponseTimeHTTP"
	val, err := parseResponse(metricName, metadata.AttributeResponseTypeHttp.String(), alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverResponseDurationDataPoint(now, val, metadata.AttributeResponseTypeHttp)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverRequestCountDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "StatNoOfRequests"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverRequestCountDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverRequestTimeoutCountDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "StatNoOfTimeouts"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverRequestTimeoutCountDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverConnectionErrorsDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "StatNoOfConnectionErrors"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverConnectionErrorsDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverCacheHitsDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "CacheHits"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverCacheHitsDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverCacheEvictionsDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "EvictedEntries"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverCacheEvictionsDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverIcmAvailabilityDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "ICM"
	val, ok := alertTreeResponse[metricName]
	if !ok {
		errs.AddPartial(1, fmt.Errorf(collectMetricError, metricName, errValueNotFound))
		return
	}

	switch models.StateColor(val) {
	case models.StateColorGray:
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 1, metadata.AttributeControlStateGrey)
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 0, metadata.AttributeControlStateGreen)
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 0, metadata.AttributeControlStateYellow)
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 0, metadata.AttributeControlStateRed)
	case models.StateColorGreen:
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 0, metadata.AttributeControlStateGrey)
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 1, metadata.AttributeControlStateGreen)
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 0, metadata.AttributeControlStateYellow)
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 0, metadata.AttributeControlStateRed)
	case models.StateColorYellow:
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 0, metadata.AttributeControlStateGrey)
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 0, metadata.AttributeControlStateGreen)
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 1, metadata.AttributeControlStateYellow)
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 0, metadata.AttributeControlStateRed)
	case models.StateColorRed:
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 0, metadata.AttributeControlStateGrey)
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 0, metadata.AttributeControlStateGreen)
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 0, metadata.AttributeControlStateYellow)
		s.mb.RecordSapnetweaverIcmAvailabilityDataPoint(now, 1, metadata.AttributeControlStateRed)
	default:
		errs.AddPartial(1, fmt.Errorf(collectMetricError, metricName, errInvalidStateColor))
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverHostCPUUtilizationDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "CPU_Utilization"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverHostCPUUtilizationDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverHostSpoolListUsedDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "HostspoolListUsed"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverHostSpoolListUsedDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func (s *sapNetweaverScraper) recordSapnetweaverShortDumpsCountDataPoint(now pcommon.Timestamp, alertTreeResponse map[string]string, errs *scrapererror.ScrapeErrors) {
	metricName := "Shortdumps Frequency"
	val, err := parseResponse(metricName, "", alertTreeResponse)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}

	err = s.mb.RecordSapnetweaverShortDumpsRateDataPoint(now, val)
	if err != nil {
		errs.AddPartial(1, err)
		return
	}
}

func parseResponse(metricName string, attributeName string, alertTreeResponse map[string]string) (string, error) {
	val, ok := alertTreeResponse[metricName]
	if !ok {
		return "", formatErrorMsg(metricName, attributeName, errValueNotFound)
	}

	if strings.Contains(val, "-") {
		return "", formatErrorMsg(metricName, attributeName, errValueHyphen)
	}
	return val, nil
}

func formatErrorMsg(metricName string, attributeName string, errValue error) error {
	if attributeName == "" {
		return fmt.Errorf(collectMetricError, metricName, errValue)
	}

	return fmt.Errorf(collectMetricErrorWithAttribute, metricName, attributeName, errValue)
}
