package logrusutil_test

import (
	"testing"

	"github.com/opalmer/logrusutil"
)

func TestNewConfig(t *testing.T) {
	cfg := logrusutil.NewConfig()
	if cfg.Level != logrusutil.DefaultLevel {
		t.Error()
	}
	if cfg.HookLevel != logrusutil.DefaultHookLevel {
		t.Error()
	}
	if cfg.HookStackLevel != logrusutil.DefaultHookStackLevel {
		t.Error()
	}
	if cfg.CallerHookField != logrusutil.DefaultCallerHookField {
		t.Error()
	}
}
