package logrusutil

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMockFieldLogger_WithFields(t *testing.T) {
	called := false
	mock := &MockFieldLogger{WithFieldsFunc: func(fields logrus.Fields) *logrus.Entry {
		called = true
		value, ok := fields["test"]
		if !ok {
			t.Error()
		}
		if value != true {
			t.Error()
		}
		return &logrus.Entry{}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.WithFields(logrus.Fields{"test": true})
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_WithField(t *testing.T) {
	called := false
	mock := &MockFieldLogger{WithFieldFunc: func(field string, value interface{}) *logrus.Entry {
		called = true
		if field != "foo" {
			t.Error()
		}
		if value != true {
			t.Error()
		}
		return &logrus.Entry{}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.WithField("foo", true)
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_WithError(t *testing.T) {
	called := false
	mock := &MockFieldLogger{WithErrorFunc: func(err error) *logrus.Entry {
		called = true
		if err == nil || err.Error() != "test" {
			t.Error()
		}
		return &logrus.Entry{}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.WithError(errors.New("test"))
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Debugf(t *testing.T) {
	called := false
	expectedFunction := Function(Debugf)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != format {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Debugf(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Infof(t *testing.T) {
	called := false
	expectedFunction := Function(Infof)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != format {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Infof(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Printf(t *testing.T) {
	called := false
	expectedFunction := Function(Printf)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != format {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Printf(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Warnf(t *testing.T) {
	called := false
	expectedFunction := Function(Warnf)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != format {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Warnf(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Warningf(t *testing.T) {
	called := false
	expectedFunction := Function(Warningf)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != format {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Warningf(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Errorf(t *testing.T) {
	called := false
	expectedFunction := Function(Errorf)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != format {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Errorf(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Fatalf(t *testing.T) {
	called := false
	expectedFunction := Function(Fatalf)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != format {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Fatalf(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Panicf(t *testing.T) {
	called := false
	expectedFunction := Function(Panicf)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != format {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Panicf(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_PanicLn(t *testing.T) {
	called := false
	expectedFunction := Function(Panicf)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != format {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Panicf(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}
func TestMockFieldLogger_Debug(t *testing.T) {
	called := false
	expectedFunction := Function(Debug)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != EmptyFormat {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Debug(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Info(t *testing.T) {
	called := false
	expectedFunction := Function(Info)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != EmptyFormat {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Info(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Print(t *testing.T) {
	called := false
	expectedFunction := Function(Print)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != EmptyFormat {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Print(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Warn(t *testing.T) {
	called := false
	expectedFunction := Function(Warn)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != EmptyFormat {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Warn(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Warning(t *testing.T) {
	called := false
	expectedFunction := Function(Warning)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != EmptyFormat {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Warning(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Error(t *testing.T) {
	called := false
	expectedFunction := Function(Error)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != EmptyFormat {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Error(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Fatal(t *testing.T) {
	called := false
	expectedFunction := Function(Fatal)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != EmptyFormat {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Fatal(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}

func TestMockFieldLogger_Panic(t *testing.T) {
	called := false
	expectedFunction := Function(Panic)
	format := "%s %s"

	mock := &MockFieldLogger{LevelFunc: func(function Function, s string, i ...interface{}) {
		called = true
		if function != expectedFunction {
			t.Error()
		}
		if s != EmptyFormat {
			t.Error()
		}
	}}

	var call = func(logger logrus.FieldLogger) {
		logger.Panic(format, "a", "b")
	}
	call(mock)
	if !called {
		t.Error()
	}
}
