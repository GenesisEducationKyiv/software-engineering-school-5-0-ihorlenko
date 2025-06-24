package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ihorlenko/weather_notifier/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BasicEndpointsTestSuite struct {
	suite.Suite
	testSuite *helpers.TestSuite
}

func (s *BasicEndpointsTestSuite) SetupSuite() {
	s.testSuite = helpers.SetupTestSuite(s.T())
}

func (s *BasicEndpointsTestSuite) TearDownSuite() {
	s.testSuite.TeardownTestSuite(s.T())
}

func (s *BasicEndpointsTestSuite) SetupTest() {
	s.testSuite.CleanupDatabase(s.T())
}

func (s *BasicEndpointsTestSuite) TestPingEndpoint() {
	resp, err := http.Get(s.testSuite.GetBaseURL() + "/ping")

	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), "pong", result["message"])
}

func (s *BasicEndpointsTestSuite) TestHealthCheck() {
	resp, err := http.Get(s.testSuite.GetBaseURL() + "/ping")
	assert.NoError(s.T(), err)
	defer resp.Body.Close()
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
}

func (s *BasicEndpointsTestSuite) TestInvalidEndpoint() {
	resp, err := http.Get(s.testSuite.GetBaseURL() + "/nonexistent")
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)
}

func TestBasicEndpointsTestSuite(t *testing.T) {
	suite.Run(t, new(BasicEndpointsTestSuite))
}
