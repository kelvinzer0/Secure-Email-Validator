package config

// Config holds the configuration for the email checker
type Config struct {
	Timeout int  // SMTP connection timeout in seconds
	Verbose bool // Enable verbose logging
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Timeout: 10,
		Verbose: false,
	}
}
