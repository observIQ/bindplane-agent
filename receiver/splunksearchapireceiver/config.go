package splunksearchapireceiver

import (
	"errors"
	"strings"
	"time"
)

type Config struct {
	Server   string   `mapstructure:"splunk_server"`
	Username string   `mapstructure:"splunk_username"`
	Password string   `mapstructure:"splunk_password"`
	Searches []Search `mapstructure:"searches"`
}

type Search struct {
	Query        string `mapstructure:"query"`
	EarliestTime string `mapstructure:"earliest_time"`
	LatestTime   string `mapstructure:"latest_time"`
	Limit        int    `mapstructure:"limit"`
}

func (cfg *Config) Validate() error {
	if cfg.Server == "" {
		return errors.New("missing Splunk server")
	}
	if cfg.Username == "" {
		return errors.New("missing Splunk username")
	}
	if cfg.Password == "" {
		return errors.New("missing Splunk password")
	}
	if len(cfg.Searches) == 0 {
		return errors.New("at least one search must be provided")
	}

	for _, search := range cfg.Searches {
		if search.Query == "" {
			return errors.New("missing query in search")
		}

		// query implicitly starts with "search" command
		if !strings.HasPrefix(search.Query, "search ") {
			search.Query = "search " + search.Query
		}

		if strings.Contains(search.Query, "|") {
			return errors.New("command chaining is not supported for queries")
		}

		if search.EarliestTime == "" {
			return errors.New("missing earliest_time in search")
		}
		if search.LatestTime == "" {
			return errors.New("missing latest_time in search")
		}

		// parse time strings to time.Time
		earliestTime, err := time.Parse(time.RFC3339, search.EarliestTime)
		if err != nil {
			return errors.New("earliest_time failed to be parsed as RFC3339")
		}

		latestTime, err := time.Parse(time.RFC3339, search.LatestTime)
		if err != nil {
			return errors.New("latest_time failed to be parsed as RFC3339")
		}

		if earliestTime.UTC().After(latestTime.UTC()) {
			return errors.New("earliest_time must be earlier than latest_time")
		}
		if earliestTime.UTC().After(time.Now().UTC()) {
			return errors.New("earliest_time must be earlier than current time")
		}
		if latestTime.UTC().After(time.Now().UTC()) {
			return errors.New("latest_time must be earlier than current time")
		}
	}
	return nil
}
