package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ihorlenko/weather_notifier/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SimpleWeatherEndpointsTestSuite tests weather endpoints without external API mocking
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

// TestGetWeatherMissingCity tests weather endpoint without city parameter
func (s *SimpleWeatherEndpointsTestSuite) TestGetWeatherMissingCity() {
	resp, err := http.Get(s.testSuite.GetBaseURL() + "/api/weather")
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Contains(s.T(), result["error"], "City parameter is required")
}

// TestGetWeatherEmptyCity tests weather endpoint with empty city parameter
func (s *SimpleWeatherEndpointsTestSuite) TestGetWeatherEmptyCity() {
	resp, err := http.Get(s.testSuite.GetBaseURL() + "/api/weather?city=")
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Contains(s.T(), result["error"], "City parameter is required")
}

// TestGetWeatherWithInvalidCity tests handling of weather API errors
func (s *SimpleWeatherEndpointsTestSuite) TestGetWeatherWithInvalidCity() {
	// This will actually call the external API, so we expect an error due to invalid API key
	resp, err := http.Get(s.testSuite.GetBaseURL() + "/api/weather?city=InvalidCity123XYZ")
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	// Should return 500 due to invalid API key or invalid city
	assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	// Should contain error about weather service
	assert.Contains(s.T(), result["error"], "weather service")
}

func TestSimpleWeatherEndpointsTestSuite(t *testing.T) {
	suite.Run(t, new(SimpleWeatherEndpointsTestSuite))
}
