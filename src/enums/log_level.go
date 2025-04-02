package enums

// LogLevel is a custom type for log levels.
type LogLevel string

// Define constants for the log levels.
const (
	Debug   LogLevel = "Debug"
	Info    LogLevel = "Info"
	Success LogLevel = "Success"
	Warning LogLevel = "Warning"
	Error   LogLevel = "Error"
	Panic   LogLevel = "panic"
)

func (l LogLevel) String() string {
	return string(l)
}
