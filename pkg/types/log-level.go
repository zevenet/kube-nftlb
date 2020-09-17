package types

// LogLevel is used to check if a message must be shown.
type LogLevel uint

const (
	// StandardLog shows only successful operations.
	StandardLog LogLevel = iota

	// ErrorLog shows errors + StandardLog.
	ErrorLog

	// DetailedLog shows every parsed struct + ErrorLog.
	DetailedLog
)
