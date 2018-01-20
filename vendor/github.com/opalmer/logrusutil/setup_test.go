package logrusutil_test

import (
	"io/ioutil"
	"testing"

	"github.com/opalmer/logrusutil"
	"github.com/sirupsen/logrus"
)

func TestConfigureRoot_ErrLevelNotProvided_Level(t *testing.T) {
	logger := logrus.New()
	cfg := logrusutil.NewConfig()
	cfg.Level = ""
	if logrusutil.ConfigureLogger(logger, cfg) != logrusutil.ErrLevelNotProvided {
		t.Error()
	}
}

func TestConfigureRoot_ErrLevelNotProvided_HookLevel(t *testing.T) {
	logger := logrus.New()
	cfg := logrusutil.NewConfig()
	cfg.HookLevel = ""
	if logrusutil.ConfigureLogger(logger, cfg) != logrusutil.ErrLevelNotProvided {
		t.Error()
	}
}

func TestConfigureRoot_DisabledLevel(t *testing.T) {
	logger := logrus.New()
	cfg := logrusutil.NewConfig()
	cfg.Level = logrusutil.DisabledLevel
	if err := logrusutil.ConfigureLogger(logger, cfg); err != nil {
		t.Error(err)
	}

	if logger.Level != logrus.PanicLevel {
		t.Error()
	}

	if logger.Out != ioutil.Discard {
		t.Error()
	}
}

func TestConfigureRoot_HookDisabled(t *testing.T) {
	logger := logrus.New()
	cfg := logrusutil.NewConfig()
	cfg.Level = logrusutil.DisabledLevel
	cfg.HookLevel = logrusutil.DisabledLevel
	if err := logrusutil.ConfigureLogger(logger, cfg); err != nil {
		t.Error(err)
	}

	if logger.Level != logrus.PanicLevel {
		t.Error()
	}

	if logger.Out != ioutil.Discard {
		t.Error()
	}
}

func TestConfigureRoot_Level_ParseError(t *testing.T) {
	logger := logrus.New()
	cfg := logrusutil.NewConfig()
	cfg.Level = "foobar"
	if err := logrusutil.ConfigureLogger(logger, cfg); err == nil {
		t.Error()
	}
}

func TestConfigureRoot_HookLevel_BadLevel(t *testing.T) {
	logger := logrus.New()
	cfg := logrusutil.NewConfig()
	cfg.HookLevel = "foobar"
	if err := logrusutil.ConfigureLogger(logger, cfg); err == nil {
		t.Error()
	}
}

func TestConfigureRoot_SetLevel(t *testing.T) {
	logger := logrus.New()
	for _, level := range logrus.AllLevels {
		cfg := logrusutil.NewConfig()
		cfg.Level = level.String()
		if err := logrusutil.ConfigureLogger(logger, cfg); err != nil {
			t.Error(err)
		}
		if logger.Level != level {
			t.Errorf("%s != %s", logger.Level, level)
		}
	}
}

func TestConfigureLogger(t *testing.T) {
	logger := logrus.New()
	cfg := logrusutil.NewConfig()
	cfg.Level = "warning"
	cfg.HookLevel = "warning"

	for _, hooks := range logger.Hooks {
		if len(hooks) != 0 {
			t.Error()
		}
	}

	if err := logrusutil.ConfigureLogger(logger, cfg); err != nil {
		t.Error(err)
	}

	for _, hooks := range logger.Hooks {
		if len(hooks) != 1 {
			t.Error()
		}
	}
}
