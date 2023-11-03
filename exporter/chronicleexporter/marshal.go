package chronicleexporter

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
)

// marshaler is an interface for marshalling logs.
//
//go:generate mockery --name logMarshaler --output ./internal/mocks --with-expecter --filename mock_marshaler.go --structname MockMarshaler
type logMarshaler interface {
	MarshalRawLogs(ld plog.Logs) ([]byte, error)
}

type marshaler struct {
	cfg Config
}

func newMarshaler(cfg Config) *marshaler {
	return &marshaler{
		cfg: cfg,
	}
}

func (ce *marshaler) MarshalRawLogs(ld plog.Logs) ([]byte, error) {
	if ce.cfg.RawLogField == "" {
		plogMarshaller := &plog.JSONMarshaler{}
		return plogMarshaller.MarshalLogs(ld)
	}

	rawLogs, err := ce.extractRawLogs(ld)
	if err != nil {
		return nil, fmt.Errorf("extract raw logs: %w", err)
	}

	rawLogData := map[string]interface{}{
		"entries":  rawLogs,
		"log_type": ce.cfg.LogType,
	}

	if ce.cfg.CustomerID != "" {
		rawLogData["custumer_id"] = ce.cfg.CustomerID
	}

	return json.Marshal(rawLogData)
}

type entry struct {
	LogText   string `json:"log_text"`
	Timestamp string `json:"timestamp"`
}

func (ce *marshaler) extractRawLogs(ld plog.Logs) ([]entry, error) {
	entries := []entry{}

	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		resourceLog := ld.ResourceLogs().At(i)
		for j := 0; j < resourceLog.ScopeLogs().Len(); j++ {
			scopeLog := resourceLog.ScopeLogs().At(j)
			for k := 0; k < scopeLog.LogRecords().Len(); k++ {
				logRecord := scopeLog.LogRecords().At(k)
				rawLog, err := ce.getRawField(logRecord)
				if err != nil {
					return nil, fmt.Errorf("get raw field: %w", err)
				}

				entries = append(entries, entry{
					LogText:   rawLog,
					Timestamp: logRecord.Timestamp().AsTime().Format(time.RFC3339Nano),
				})
			}
		}
	}

	return entries, nil
}

func (ce *marshaler) getRawField(logRecord plog.LogRecord) (string, error) {
	topLevelField, nestedFields, err := parseLogField(ce.cfg.RawLogField)
	if err != nil {
		return "", err
	}

	if len(nestedFields) == 0 {
		return ce.getTopLevelFieldAsString(logRecord, topLevelField)
	}

	var logMap map[string]any
	switch topLevelField {
	case "attributes":
		logMap = logRecord.Attributes().AsRaw()
	case "body":
		if logRecord.Body().Type() != pcommon.ValueTypeMap {
			return "", errors.New("body is not a map")
		}
		logMap = logRecord.Body().Map().AsRaw()
	default:
		return "", fmt.Errorf("unsupported top level field: %s", topLevelField)
	}

	return extractNestedValue(logMap, nestedFields)
}

func (ce *marshaler) getTopLevelFieldAsString(logRecord plog.LogRecord, field string) (string, error) {
	switch field {
	case "attributes":
		attributes := logRecord.Attributes().AsRaw()
		bytes, err := json.Marshal(attributes)
		if err != nil {
			return "", fmt.Errorf("failed to marshal attributes: %w", err)
		}
		return string(bytes), nil
	case "body":
		switch logRecord.Body().Type() {
		case pcommon.ValueTypeStr:
			return logRecord.Body().Str(), nil
		case pcommon.ValueTypeMap:
			bodyMap := logRecord.Body().Map().AsRaw()
			bytes, err := json.Marshal(bodyMap)
			if err != nil {
				return "", fmt.Errorf("failed to marshal body map: %w", err)
			}
			return string(bytes), nil
		default:
			return "", errors.New("unsupported body type")
		}
	default:
		return "", fmt.Errorf("unsupported top level field: %s", field)
	}
}

func parseLogField(field string) (string, []string, error) {
	parts := strings.SplitN(field, `["`, 2)

	if len(parts) == 1 {
		return parts[0], nil, nil
	}

	re := regexp.MustCompile(`\["(.*?)"\]`)
	matches := re.FindAllStringSubmatch(field, -1)

	keys := make([]string, len(matches))
	for i, match := range matches {
		if len(match) > 1 {
			keys[i] = match[1]
		}
	}

	return parts[0], keys, nil
}

func extractNestedValue(logMap map[string]any, keys []string) (string, error) {
	for i, key := range keys {
		value, ok := logMap[key]
		if !ok {
			return "", fmt.Errorf("failed to find key '%s' in log map", key)
		}

		if i == len(keys)-1 {
			if strVal, ok := value.(string); ok {
				return strVal, nil
			}
			return "", errors.New("final value is not a string")
		}

		nextMap, ok := value.(map[string]any)
		if !ok {
			return "", fmt.Errorf("value for key %s is not a map", key)
		}
		logMap = nextMap
	}

	return "", fmt.Errorf("failed to parse raw log field")
}
