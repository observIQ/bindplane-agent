package bindplaneexporter

import (
	"encoding/json"
	"strings"
)

// LiveTailConfig is a live tail config file
type LiveTailConfig struct {
	Sessions []Session `yaml:"sessions"`
}

// Session is a bindplane session
type Session struct {
	ID      string   `yaml:"id"`
	Filters []string `yaml:"filters"`
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
