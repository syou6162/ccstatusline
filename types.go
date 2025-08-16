package main

// Config represents the configuration structure
type Config struct {
	Actions   []Action `yaml:"actions"`
	Separator string   `yaml:"separator"`
}

// Action represents a single action in the configuration
type Action struct {
	Name    string  `yaml:"name"`
	Command Command `yaml:"command"`
}

// Command represents the command configuration
type Command struct {
	Type    string `yaml:"type"`    // "command" or "output"
	Command string `yaml:"command"` // for command type
	Text    string `yaml:"text"`    // for output type
	Color   string `yaml:"color"`   // for output type
}

// Context holds the processing context
type Context struct {
	InputJSON      map[string]interface{} // JSON from Claude Code
	CommandOutputs map[string]string      // Cache of command outputs
	CurrentOutput  string                 // Current command output
}
