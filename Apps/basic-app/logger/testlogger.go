package logger

import "testing"

type LoggerForTests struct {
	Tester *testing.T
}

func (l LoggerForTests) Debug(args ...interface{}) {
	l.Tester.Log(args...)
}

func (l LoggerForTests) Debugf(format string, args ...interface{}) {
	l.Tester.Logf(format, args...)
}

func (l LoggerForTests) Info(args ...interface{}) {
	l.Tester.Log(args...)
}

func (l LoggerForTests) Infof(format string, args ...interface{}) {
	l.Tester.Logf(format, args...)
}

func (l LoggerForTests) Warning(args ...interface{}) {
	l.Tester.Log(args...)
}

func (l LoggerForTests) Warningf(format string, args ...interface{}) {
	l.Tester.Logf(format, args...)
}

func (l LoggerForTests) Error(args ...interface{}) {
	l.Tester.Log(args...)
}

func (l LoggerForTests) Errorf(format string, args ...interface{}) {
	l.Tester.Logf(format, args...)
}
