package logrusutil_test

import (
	"testing"

	"github.com/opalmer/logrusutil"
	"github.com/sirupsen/logrus"
)

func TestExampleConfigureRoot(t *testing.T) {
	ExampleConfigureRoot()
}

// NOTE: This example is included in README.md. If you update it, please update
// the readme too.

func ExampleConfigureRoot() {
	// Setup the root logger and hooks.
	if err := logrusutil.ConfigureLogger(logrus.StandardLogger(), logrusutil.NewConfig()); err != nil {
		panic(err)
	}
}
