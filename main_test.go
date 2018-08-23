package xaqt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

// entry point for all package internal tests.
// tests are grouped into suites according to the code being tested.
//
func TestMain(m *testing.M) {
	retCode := m.Run()

	os.Exit(retCode)
}

func TestSandboxSuite(t *testing.T) {
	suite.Run(t, &SandboxTestSuite{})
}
