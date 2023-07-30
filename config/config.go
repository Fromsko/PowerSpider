package config

// Config represents the configuration with various fields.
type Config struct {
	User    string // Username for authentication.
	Pwd     string // Password for authentication.
	Timer   Timer  // Program run interval.
	ResDir  string // Directory to save the final results.
	Porxy   string // Proxy for accessing the school network from outside.
	BaseUrl string // Base URL for the HomePage.
}

// Timer represents the time information.
type Timer struct {
	TimeUnit string // Time unit.
	TimeInfo int    // Time information.
}

func InitConfig(c *Config) *Config {
	// replaceIfEmpty is a function to check if a field is empty and replace it with a default value.
	replaceIfEmpty := func(field *string, defaultValue string) {
		if *field == "" {
			*field = defaultValue
		}
	}

	if c != nil {
		replaceIfEmpty(&c.ResDir, "res")
		replaceIfEmpty(&c.BaseUrl, HomeUrl)
		return c
	}

	c = &Config{
		ResDir: "res",
		Timer: Timer{
			TimeUnit: "hours",
			TimeInfo: 1,
		},
		BaseUrl: HomeUrl,
	}

	return c
}

// SetProxy sets the Porxy field in the Config object to the specified URL.
func (c *Config) SetProxy(url string) {
	c.Porxy = url
}

// SetBaseUrl sets the BaseUrl field in the Config object to the specified URL.
func (c *Config) SetBaseUrl(url string) {
	c.BaseUrl = url
}

// SetTimer sets the Timer field in the Config object to the specified Timer value.
func (c *Config) SetTimer(t Timer) {
	c.Timer = t
}
