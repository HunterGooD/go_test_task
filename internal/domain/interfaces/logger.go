package interfaces

type Logger interface {
	Info(message string, opt ...any)
	Debug(message string, opt ...any)
	Warn(message string, opt ...any)
	Error(message string, opt ...any)
}
