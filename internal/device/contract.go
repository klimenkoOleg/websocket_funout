package device

type Logger interface {
	Debug(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Panic(args ...interface{})
}
