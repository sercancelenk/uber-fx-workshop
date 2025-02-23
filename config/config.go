package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type RootConfig struct {
	KafkaAthena KafkaAthenaConfig `mapstructure:"config"`
}

// Structs for Configuration
type KafkaAthenaConfig struct {
	Clusters  map[string]ClusterConfig  `mapstructure:"clusters"`
	Consumers map[string]ConsumerConfig `mapstructure:"consumers"`
}

type ClusterConfig struct {
	Servers []string          `mapstructure:"servers"`
	Props   map[string]string `mapstructure:"props"`
}

type ConsumerConfig struct {
	Cluster  string   `mapstructure:"cluster"`
	Topic    string   `mapstructure:"topic"`
	Failover Failover `mapstructure:"failover"`
}

type Failover struct {
	DLQ DLQConfig `mapstructure:"dlq"`
}

type DLQConfig struct {
	Topic string `mapstructure:"topic"`
}

// Mutex for thread safety
var mu sync.Mutex

// LoadConfig initializes and returns the parsed configuration
func LoadConfig() (*RootConfig, error) {
	viper.SetConfigName("kafkathena")                  // Config file name (without extension)
	viper.SetConfigType("yaml")                        // Explicitly set file type
	viper.AddConfigPath(os.Getenv("CONFIG_FILE_PATH")) // Load from environment variable
	viper.AddConfigPath(".")                           // Also search current directory

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	config := &RootConfig{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Watch for changes and refresh configuration
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		mu.Lock()
		defer mu.Unlock()
		if err := viper.Unmarshal(config); err != nil {
			zap.L().Error("error when refreshing the config", zap.Error(err))
		}
	})

	return config, nil
}
