package internal

type Logger interface {
	Debug(message string)
	Debugf(message string, args ...interface{})
	Info(message string)
	Infof(message string, args ...interface{})
	Warn(message string)
	Warnf(message string, args ...interface{})
	Fatal(err error)
	Fatalf(message string, args ...interface{})
}
