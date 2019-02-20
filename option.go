package logrus

import "github.com/sirupsen/logrus"

// Option is a function modifying the logger and used as option.
type Option = func(*logrus.Logger)

// WithFormatter update the logger format.
func WithFormatter(formatter logrus.Formatter) Option {
	return func(logger *logrus.Logger) {
		logger.Formatter = formatter
	}
}
