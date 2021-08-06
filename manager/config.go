package manager

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	endpoint          = "wss://connections.app.observiq.com"
	statusInterval    = time.Minute
	reconnectInterval = time.Minute * 30
	maxConnectBackoff = time.Minute * 5
	bufferSize        = 50
)

// Config is the configuration of an observiq extension
type Config struct {
	Endpoint          string        `mapstructure:"endpoint"`
	AgentName         string        `mapstructure:"agent_name"`
	AgentID           string        `mapstructure:"agent_id"`
	SecretKey         string        `mapstructure:"secret_key"`
	StatusInterval    time.Duration `mapstructure:"status_interval"`
	ReconnectInterval time.Duration `mapstructure:"reconnect_interval"`
	MaxConnectBackoff time.Duration `mapstructure:"max_connect_backoff"`
	BufferSize        int           `mapstructure:"buffer_size"`
	TemplateID        string        `mapstructure:"template_id"`
}

// ConfigFromFile creates a config from the supplied file
func ConfigFromFile(filePath string) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetDefault("endpoint", endpoint)
	viper.SetDefault("status_interval", statusInterval)
	viper.SetDefault("reconnect_interval", reconnectInterval)
	viper.SetDefault("max_connect_backoff", maxConnectBackoff)
	viper.SetDefault("buffer_size", bufferSize)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if err := viper.ReadConfig(file); err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return config, nil
}
