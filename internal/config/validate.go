package config

import "errors"

var (
	ErrInvalidLogLevel = errors.New("invalid log level provided")
	ErrROMNotProvided  = errors.New("ROM file path must be provided")
)

func (c Config) Validate() error {
	if c.LogLevel != LogLevelDebug &&
		c.LogLevel != LogLevelInfo &&
		c.LogLevel != LogLevelWarn &&
		c.LogLevel != LogLevelError {
		return ErrInvalidLogLevel
	}

	if c.ROM == "" {
		return ErrROMNotProvided
	}

	return nil
}
