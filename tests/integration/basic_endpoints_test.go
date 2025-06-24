package integration

import (
	"context"
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
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.testSuite.GetBaseURL()+"/ping", nil)
	assert.NoError(s.T(), err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), "pong", result["message"])
}

func (s *BasicEndpointsTestSuite) TestHealthCheck() {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.testSuite.GetBaseURL()+"/ping", nil)
	assert.NoError(s.T(), err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)
}

func (s *BasicEndpointsTestSuite) TestInvalidEndpoint() {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.testSuite.GetBaseURL()+"/nonexistent", nil)
	assert.NoError(s.T(), err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)
}

func TestBasicEndpointsTestSuite(t *testing.T) {
	suite.Run(t, new(BasicEndpointsTestSuite))
}
