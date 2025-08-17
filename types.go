package main

// Config represents the configuration structure
type Config struct {
	Actions   []Action `yaml:"actions"`
	Separator string   `yaml:"separator"`
}

// Action represents a single action in the configuration
type Action struct {
	Name    string `yaml:"name"`    // Optional name for debugging
	Command string `yaml:"command"` // Shell command to execute or template text
	Color   string `yaml:"color"`   // Optional color
}
