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

type SimpleWeatherEndpointsTestSuite struct {
	suite.Suite
	testSuite *helpers.TestSuite
}

func (s *SimpleWeatherEndpointsTestSuite) SetupSuite() {
	s.testSuite = helpers.SetupTestSuite(s.T())
}

func (s *SimpleWeatherEndpointsTestSuite) TearDownSuite() {
	s.testSuite.TeardownTestSuite(s.T())
}

func (s *SimpleWeatherEndpointsTestSuite) SetupTest() {
	s.testSuite.CleanupDatabase(s.T())
}

func (s *SimpleWeatherEndpointsTestSuite) TestGetWeatherMissingCity() {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.testSuite.GetBaseURL()+"/api/weather", nil)
	assert.NoError(s.T(), err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Contains(s.T(), result["error"], "City parameter is required")
}

func (s *SimpleWeatherEndpointsTestSuite) TestGetWeatherEmptyCity() {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.testSuite.GetBaseURL()+"/api/weather?city=", nil)
	assert.NoError(s.T(), err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Contains(s.T(), result["error"], "City parameter is required")
}

func (s *SimpleWeatherEndpointsTestSuite) TestGetWeatherWithInvalidCity() {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet,
		s.testSuite.GetBaseURL()+"/api/weather?city=InvalidCity123XYZ", nil,
	)
	assert.NoError(s.T(), err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Contains(s.T(), result["error"], "weather service")
}

func TestSimpleWeatherEndpointsTestSuite(t *testing.T) {
	suite.Run(t, new(SimpleWeatherEndpointsTestSuite))
}
