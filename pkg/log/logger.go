package log

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, err error, keysAndValues ...interface{})
	Error(msg string, err error, keysAndValues ...interface{})
	Fatal(msg string, err error, keysAndValues ...interface{})
	With(keysAndValues ...interface{}) Logger
}
