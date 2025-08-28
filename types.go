package main

// Config represents the configuration structure
type Config struct {
	Actions   []Action `yaml:"actions"`
	Separator string   `yaml:"separator"`
}

// Action represents a single action in the configuration
type Action struct {
	Name     string `yaml:"name"`      // Required: unique identifier for action
	Type     string `yaml:"type"`      // "command" (default) or "builtin"
	Command  string `yaml:"command"`   // Shell command to execute (used for command type)
	Function string `yaml:"function"`  // Builtin function name (used for builtin type)
	Prefix   string `yaml:"prefix"`    // Optional prefix to prepend to command output
	Color    string `yaml:"color"`     // Optional color (foreground or background with bg_ prefix)
	CacheTTL int    `yaml:"cache_ttl"` // Cache TTL in seconds (0 or unset = no cache)
}
