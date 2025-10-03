package log

type Option func(*Logger)

func WithMessageKey(key string) Option {
	return func(logger *Logger) {
		logger.msgKey = key
	}
}

func WithSprint(sprint func(...any) string) Option {
	return func(logger *Logger) {
		logger.sprint = sprint
	}
}

func WithSprintf(sprintf func(format string, a ...any) string) Option {
	return func(logger *Logger) {
		logger.sprintf = sprintf
	}
}
