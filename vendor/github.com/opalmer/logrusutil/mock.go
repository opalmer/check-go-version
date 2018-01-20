package logrusutil

import "github.com/sirupsen/logrus"

// Function is a type used by the constants below.
type Function uint32

const (
	// NoLevel and the following constants are used by MockFieldLogger to
	// indicate to MockFieldLogger.levelFunc what function is being called by
	// logrus.
	NoLevel = iota
	Debugf  // nolint
	Infof
	Printf
	Warnf
	Warningf
	Errorf
	Fatalf
	Panicf
	Debug
	Info
	Print
	Warn
	Warning
	Error
	Fatal
	Panic
	Debugln
	Infoln
	Println
	Warnln
	Warningln
	Errorln
	Fatalln
	Panicln

	// EmptyFormat is passed to MockFieldLogger.LevelFunc when there's
	// no format associated with a log entry.
	EmptyFormat = ""
)

// MockFieldLogger implements logrus.FieldLogger. This is intended for use
// in tests.
type MockFieldLogger struct {
	WithFieldFunc  func(string, interface{}) *logrus.Entry
	WithFieldsFunc func(logrus.Fields) *logrus.Entry
	WithErrorFunc  func(error) *logrus.Entry

	// LevelFunc is intended to provide a single implementation point
	// for all logging level both formatted and unformatted. You may wish
	// to review the implementation of this struct to understand how this
	// function is called in various cases.
	LevelFunc func(Function, string, ...interface{})
}

// WithField ...
func (m *MockFieldLogger) WithField(key string, value interface{}) *logrus.Entry {
	return m.WithFieldFunc(key, value)
}

// WithFields ...
func (m *MockFieldLogger) WithFields(fields logrus.Fields) *logrus.Entry {
	return m.WithFieldsFunc(fields)
}

// WithError ...
func (m *MockFieldLogger) WithError(err error) *logrus.Entry {
	return m.WithErrorFunc(err)
}

// Debugf ...
func (m *MockFieldLogger) Debugf(format string, args ...interface{}) {
	m.LevelFunc(Debugf, format, args)
}

// Infof ...
func (m *MockFieldLogger) Infof(format string, args ...interface{}) {
	m.LevelFunc(Infof, format, args)
}

// Printf ...
func (m *MockFieldLogger) Printf(format string, args ...interface{}) {
	m.LevelFunc(Printf, format, args)
}

// Warnf ...
func (m *MockFieldLogger) Warnf(format string, args ...interface{}) {
	m.LevelFunc(Warnf, format, args)
}

// Warningf ...
func (m *MockFieldLogger) Warningf(format string, args ...interface{}) {
	m.LevelFunc(Warningf, format, args)
}

// Errorf ...
func (m *MockFieldLogger) Errorf(format string, args ...interface{}) {
	m.LevelFunc(Errorf, format, args)
}

// Fatalf ...
func (m *MockFieldLogger) Fatalf(format string, args ...interface{}) {
	m.LevelFunc(Fatalf, format, args)
}

// Panicf ...
func (m *MockFieldLogger) Panicf(format string, args ...interface{}) {
	m.LevelFunc(Panicf, format, args)
}

// Debug ...
func (m *MockFieldLogger) Debug(args ...interface{}) {
	m.LevelFunc(Debug, EmptyFormat, args)
}

// Info ...
func (m *MockFieldLogger) Info(args ...interface{}) {
	m.LevelFunc(Info, EmptyFormat, args)
}

// Print ...
func (m *MockFieldLogger) Print(args ...interface{}) {
	m.LevelFunc(Print, EmptyFormat, args)
}

// Warn ...
func (m *MockFieldLogger) Warn(args ...interface{}) {
	m.LevelFunc(Warn, EmptyFormat, args)
}

// Warning ...
func (m *MockFieldLogger) Warning(args ...interface{}) {
	m.LevelFunc(Warning, EmptyFormat, args)
}

// Error ...
func (m *MockFieldLogger) Error(args ...interface{}) {
	m.LevelFunc(Error, EmptyFormat, args)
}

// Fatal ...
func (m *MockFieldLogger) Fatal(args ...interface{}) {
	m.LevelFunc(Fatal, EmptyFormat, args)
}

// Panic ...
func (m *MockFieldLogger) Panic(args ...interface{}) {
	m.LevelFunc(Panic, EmptyFormat, args)
}

// Debugln ...
func (m *MockFieldLogger) Debugln(args ...interface{}) {
	m.LevelFunc(Debugln, EmptyFormat, args)
}

// Infoln ...
func (m *MockFieldLogger) Infoln(args ...interface{}) {
	m.LevelFunc(Infoln, EmptyFormat, args)
}

// Println ...
func (m *MockFieldLogger) Println(args ...interface{}) {
	m.LevelFunc(Println, EmptyFormat, args)
}

// Warnln ...
func (m *MockFieldLogger) Warnln(args ...interface{}) {
	m.LevelFunc(Warnln, EmptyFormat, args)
}

// Warningln ...
func (m *MockFieldLogger) Warningln(args ...interface{}) {
	m.LevelFunc(Warningln, EmptyFormat, args)
}

// Errorln ...
func (m *MockFieldLogger) Errorln(args ...interface{}) {
	m.LevelFunc(Errorln, EmptyFormat, args)
}

// Fatalln ...
func (m *MockFieldLogger) Fatalln(args ...interface{}) {
	m.LevelFunc(Fatalln, EmptyFormat, args)
}

// Panicln ...
func (m *MockFieldLogger) Panicln(args ...interface{}) {
	m.LevelFunc(Panicln, EmptyFormat, args)
}
