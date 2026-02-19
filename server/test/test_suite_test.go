package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestSuite runs all the tests
func TestAccountServerTestSuite(t *testing.T) {
	suite.Run(t, new(AccountServerTestSuite))
}

// AccountServerTestSuite is the main test suite
type AccountServerTestSuite struct {
	suite.Suite
}

func (suite *AccountServerTestSuite) SetupSuite() {
	// Setup code before all tests
}

func (suite *AccountServerTestSuite) TearDownSuite() {
	// Cleanup code after all tests
}

func (suite *AccountServerTestSuite) SetupTest() {
	// Setup code before each test
}

func (suite *AccountServerTestSuite) TearDownTest() {
	// Cleanup code after each test
}
