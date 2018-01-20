package logrusutil_test

import (
	"testing"

	"github.com/opalmer/logrusutil"
	"github.com/sirupsen/logrus"
)

func TestCallerHook_Fire_Disabled(t *testing.T) {
	entry := &logrus.Entry{}
	hook := logrusutil.NewCallerHook(
		true, logrusutil.DefaultHookStackLevel,
		"test", logrus.DebugLevel)

	if err := hook.Fire(entry); err != nil {
		t.Error()
	}

	if _, ok := entry.Data["test"]; ok {
		t.Error()
	}
}

func TestCallerHook_Fire(t *testing.T) {
	entry := logrus.WithField("test", "foo")
	hook := logrusutil.NewCallerHook(
		false, logrusutil.DefaultHookStackLevel,
		"test", logrus.DebugLevel)

	if err := hook.Fire(entry); err != nil {
		t.Error()
	}

	value, ok := entry.Data["test"]
	if !ok {
		t.Error("Missing 'test' field")
	}

	if value == "foo" {
		t.Error()
	}
}

func TestCallerHook_Fire_NilMapDoesNotPanic(t *testing.T) {
	entry := &logrus.Entry{}
	hook := logrusutil.NewCallerHook(
		false, logrusutil.DefaultHookStackLevel,
		"test", logrus.DebugLevel)
	if err := hook.Fire(entry); err != nil { // Should not panic
		t.Error()
	}
}

func TestCallerHook_Levels(t *testing.T) {
	hook := logrusutil.NewCallerHook(
		false, logrusutil.DefaultHookStackLevel,
		"test", logrus.WarnLevel)

	expected := []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
	levels := hook.Levels()
	if len(levels) != len(expected) {
		t.Errorf("%d != %d", len(levels), len(expected))
	}

	for i, lvl := range expected {
		if levels[i] != lvl {
			t.Errorf("%s != %s", levels[i], lvl)
		}
	}
}

func TestCallerHook_LevelsDisabled(t *testing.T) {
	hook := logrusutil.NewCallerHook(
		true, logrusutil.DefaultHookStackLevel,
		"test", logrus.WarnLevel)
	if len(hook.Levels()) != 0 {
		t.Error()
	}
}
