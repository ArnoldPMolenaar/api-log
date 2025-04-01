package enums

// LogLevel is a custom type for log levels.
type LogLevel string

// Define constants for the log levels.
const (
	Debug   LogLevel = "Debug"
	Info             = "Info"
	Success          = "Success"
	Warning          = "Warning"
	Error            = "Error"
	Panic            = "panic"
)
