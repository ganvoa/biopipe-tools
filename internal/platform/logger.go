package platform

import "github.com/sirupsen/logrus"

type logger struct {
	tag    string
	logger *logrus.Logger
}

func NewLogger(tag string, debug bool) *logger {
	l := &logger{}
	l.tag = tag
	l.logger = logrus.New()
	l.logger.SetFormatter(&logrus.TextFormatter{})
	if debug {
		l.logger.Level = logrus.DebugLevel
	}
	return l
}

func (l *logger) Debugf(message string, args ...interface{}) {
	l.logger.WithField("tag", l.tag).Debugf(message, args)
}

func (l *logger) Debug(message string) {
	l.logger.WithField("tag", l.tag).Debug(message)
}

func (l *logger) Infof(message string, args ...interface{}) {
	l.logger.WithField("tag", l.tag).Infof(message, args)
}

func (l *logger) Info(message string) {
	l.logger.WithField("tag", l.tag).Info(message)
}

func (l *logger) Warnf(message string, args ...interface{}) {
	l.logger.WithField("tag", l.tag).Warnf(message, args)
}

func (l *logger) Warn(message string) {
	l.logger.WithField("tag", l.tag).Warn(message)
}

func (l *logger) Fatal(err error) {
	l.logger.WithError(err).WithField("tag", l.tag).Fatal("an error occured")
}

func (l *logger) Fatalf(message string, args ...interface{}) {
	l.logger.WithField("tag", l.tag).Fatalf(message, args)
}
