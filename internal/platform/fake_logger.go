package platform

type FakeLogger struct {
}

func (l *FakeLogger) Debugf(message string, args ...interface{}) {
}

func (l *FakeLogger) Debug(message string) {
}

func (l *FakeLogger) Infof(message string, args ...interface{}) {
}

func (l *FakeLogger) Info(message string) {
}

func (l *FakeLogger) Warnf(message string, args ...interface{}) {
}

func (l *FakeLogger) Warn(message string) {
}

func (l *FakeLogger) Fatal(err error) {
}

func (l *FakeLogger) Fatalf(message string, args ...interface{}) {
}
