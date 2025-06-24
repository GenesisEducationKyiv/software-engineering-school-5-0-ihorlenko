package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ihorlenko/weather_notifier/internal/models"
	"github.com/ihorlenko/weather_notifier/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SimpleSubscriptionEndpointsTestSuite tests subscription endpoints without external API mocking
type SimpleSubscriptionEndpointsTestSuite struct {
	suite.Suite
	testSuite *helpers.TestSuite
}

func (s *SimpleSubscriptionEndpointsTestSuite) SetupSuite() {
	s.testSuite = helpers.SetupTestSuite(s.T())
}

func (s *SimpleSubscriptionEndpointsTestSuite) TearDownSuite() {
	s.testSuite.TeardownTestSuite(s.T())
}

func (s *SimpleSubscriptionEndpointsTestSuite) SetupTest() {
	s.testSuite.CleanupDatabase(s.T())
}

// TestSubscribeInvalidEmail tests subscription with invalid email
func (s *SimpleSubscriptionEndpointsTestSuite) TestSubscribeInvalidEmail() {
	subscribeRequest := map[string]interface{}{
		"email":     "invalid-email",
		"city":      "London",
		"frequency": "daily",
	}

	jsonBody, err := json.Marshal(subscribeRequest)
	assert.NoError(s.T(), err)

	resp, err := http.Post(
		s.testSuite.GetBaseURL()+"/api/subscribe",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Contains(s.T(), result["error"], "Invalid request format")
}

// TestSubscribeInvalidFrequency tests subscription with invalid frequency
func (s *SimpleSubscriptionEndpointsTestSuite) TestSubscribeInvalidFrequency() {
	subscribeRequest := map[string]interface{}{
		"email":     "test@example.com",
		"city":      "London",
		"frequency": "weekly", // Invalid frequency
	}

	jsonBody, err := json.Marshal(subscribeRequest)
	assert.NoError(s.T(), err)

	resp, err := http.Post(
		s.testSuite.GetBaseURL()+"/api/subscribe",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Contains(s.T(), result["error"], "Invalid request format")
}

// TestSubscribeWithInvalidCity tests subscription with invalid city (will hit external API)
func (s *SimpleSubscriptionEndpointsTestSuite) TestSubscribeWithInvalidCity() {
	subscribeRequest := map[string]interface{}{
		"email":     "test@example.com",
		"city":      "InvalidCity123XYZ", // Very unlikely to be a real city
		"frequency": "daily",
	}

	jsonBody, err := json.Marshal(subscribeRequest)
	assert.NoError(s.T(), err)

	resp, err := http.Post(
		s.testSuite.GetBaseURL()+"/api/subscribe",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	// Should fail due to city validation (external API call will fail)
	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Contains(s.T(), result["error"], "Invalid city")
}

// TestConfirmSubscriptionInvalidToken tests confirmation with invalid token
func (s *SimpleSubscriptionEndpointsTestSuite) TestConfirmSubscriptionInvalidToken() {
	// Create HTTP client that doesn't follow redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}

	resp, err := client.Get(s.testSuite.GetBaseURL() + "/api/confirm/invalid-token")
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusFound, resp.StatusCode)

	location := resp.Header.Get("Location")
	assert.Contains(s.T(), location, "message_type=error")
	assert.Contains(s.T(), location, "not+found")
}

// TestUnsubscribeInvalidToken tests unsubscription with invalid token
func (s *SimpleSubscriptionEndpointsTestSuite) TestUnsubscribeInvalidToken() {
	// Create HTTP client that doesn't follow redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}

	resp, err := client.Get(s.testSuite.GetBaseURL() + "/api/unsubscribe/invalid-token")
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusFound, resp.StatusCode)

	location := resp.Header.Get("Location")
	assert.Contains(s.T(), location, "message_type=error")
}

// TestSubscribeMissingFields tests subscription with missing required fields
func (s *SimpleSubscriptionEndpointsTestSuite) TestSubscribeMissingFields() {
	subscribeRequest := map[string]interface{}{
		"email": "test@example.com",
		// Missing city and frequency
	}

	jsonBody, err := json.Marshal(subscribeRequest)
	assert.NoError(s.T(), err)

	resp, err := http.Post(
		s.testSuite.GetBaseURL()+"/api/subscribe",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
}

// TestConfirmUnsubscribeFlow tests the complete flow without actual email sending
func (s *SimpleSubscriptionEndpointsTestSuite) TestConfirmUnsubscribeFlow() {
	// Create HTTP client that doesn't follow redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}

	// Create a subscription directly in the database
	user := &models.User{Email: "test@example.com"}
	err := s.testSuite.DB.Create(user).Error
	assert.NoError(s.T(), err)

	subscription := &models.Subscription{
		UserID:            user.ID,
		City:              "TestCity",
		Frequency:         "daily",
		Status:            "pending",
		ConfirmationToken: "test-confirmation-token",
		UnsubscribeToken:  "test-unsubscribe-token",
	}
	err = s.testSuite.DB.Create(subscription).Error
	assert.NoError(s.T(), err)

	// Test confirmation
	resp, err := client.Get(s.testSuite.GetBaseURL() + "/api/confirm/test-confirmation-token")
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusFound, resp.StatusCode)
	location := resp.Header.Get("Location")
	assert.Contains(s.T(), location, "message_type=success")

	// Verify subscription status was updated
	var updatedSubscription models.Subscription
	err = s.testSuite.DB.Where("id = ?", subscription.ID).First(&updatedSubscription).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "active", updatedSubscription.Status)

	// Test unsubscription
	resp, err = client.Get(s.testSuite.GetBaseURL() + "/api/unsubscribe/test-unsubscribe-token")
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusFound, resp.StatusCode)
	location = resp.Header.Get("Location")
	assert.Contains(s.T(), location, "message_type=success")

	// Verify subscription status was updated to cancelled
	err = s.testSuite.DB.Where("id = ?", subscription.ID).First(&updatedSubscription).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "cancelled", updatedSubscription.Status)
}

func TestSimpleSubscriptionEndpointsTestSuite(t *testing.T) {
	suite.Run(t, new(SimpleSubscriptionEndpointsTestSuite))
}
