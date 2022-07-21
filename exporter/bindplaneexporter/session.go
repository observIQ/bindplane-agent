package bindplaneexporter

import (
	"encoding/json"
	"strings"
)

// Session is a bindplane session
type Session struct {
	ID      string   `json:"id"`
	Filters []string `json:"filters"`
}

// Matches determines if a record matches the session filters
func (s *Session) Matches(record interface{}) bool {
	if len(s.Filters) == 0 {
		return true
	}

	jsonBytes, err := json.Marshal(record)
	if err != nil {
		return false
	}
	jsonStr := string(jsonBytes)

	for _, filter := range s.Filters {
		if !strings.Contains(jsonStr, filter) {
			return false
		}
	}

	return true
}
