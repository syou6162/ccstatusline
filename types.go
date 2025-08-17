package main

// Config represents the configuration structure
type Config struct {
	Actions   []Action `yaml:"actions"`
	Separator string   `yaml:"separator"`
}

// Action represents a single action in the configuration
type Action struct {
	Name     string `yaml:"name"`      // Required: unique identifier for action
	Command  string `yaml:"command"`   // Shell command to execute or template text
	Prefix   string `yaml:"prefix"`    // Optional prefix to prepend to command output
	Color    string `yaml:"color"`     // Optional color
	CacheTTL int    `yaml:"cache_ttl"` // Cache TTL in seconds (0 or unset = no cache)
}
