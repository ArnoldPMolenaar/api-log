package enums

// LogLevel is a custom type for log levels.
type LogLevel int

// Define constants for the log levels.
const (
	Debug LogLevel = iota
	Info
	Success
	Warning
	Error
	Panic
)

// String returns the string representation of the LogLevel.
func (l LogLevel) String() string {
	return [...]string{"Debug", "Info", "Success", "Warning", "Error", "Panic"}[l]
}
