package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ihorlenko/weather_notifier/internal/models"
	"github.com/ihorlenko/weather_notifier/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

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

func (s *SimpleSubscriptionEndpointsTestSuite) testSubscribeRequest(
	subscribeRequest map[string]interface{},
	expectedErrorContent string,
) {
	jsonBody, err := json.Marshal(subscribeRequest)
	assert.NoError(s.T(), err)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		s.testSuite.GetBaseURL()+"/api/subscribe",
		bytes.NewBuffer(jsonBody))
	assert.NoError(s.T(), err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Contains(s.T(), result["error"], expectedErrorContent)
}

func (s *SimpleSubscriptionEndpointsTestSuite) TestSubscribeInvalidEmail() {
	subscribeRequest := map[string]interface{}{
		"email":     "invalid-email",
		"city":      "London",
		"frequency": "daily",
	}

	s.testSubscribeRequest(subscribeRequest, "Invalid request format")
}

func (s *SimpleSubscriptionEndpointsTestSuite) TestSubscribeInvalidFrequency() {
	subscribeRequest := map[string]interface{}{
		"email":     "test@example.com",
		"city":      "London",
		"frequency": "weekly",
	}

	s.testSubscribeRequest(subscribeRequest, "Invalid request format")
}

func (s *SimpleSubscriptionEndpointsTestSuite) TestSubscribeWithInvalidCity() {
	subscribeRequest := map[string]interface{}{
		"email":     "test@example.com",
		"city":      "InvalidCity123XYZ",
		"frequency": "daily",
	}

	s.testSubscribeRequest(subscribeRequest, "Invalid city")
}

func (s *SimpleSubscriptionEndpointsTestSuite) TestConfirmSubscriptionInvalidToken() {
	client := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		s.testSuite.GetBaseURL()+"/api/confirm/invalid-token", nil)
	assert.NoError(s.T(), err)

	resp, err := client.Do(req)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusFound, resp.StatusCode)

	location := resp.Header.Get("Location")
	assert.Contains(s.T(), location, "message_type=error")
	assert.Contains(s.T(), location, "not+found")
}

func (s *SimpleSubscriptionEndpointsTestSuite) TestUnsubscribeInvalidToken() {
	client := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		s.testSuite.GetBaseURL()+"/api/unsubscribe/invalid-token", nil)
	assert.NoError(s.T(), err)

	resp, err := client.Do(req)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusFound, resp.StatusCode)

	location := resp.Header.Get("Location")
	assert.Contains(s.T(), location, "message_type=error")
}

func (s *SimpleSubscriptionEndpointsTestSuite) TestSubscribeMissingFields() {
	subscribeRequest := map[string]interface{}{
		"email": "test@example.com",
	}

	jsonBody, err := json.Marshal(subscribeRequest)
	assert.NoError(s.T(), err)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		s.testSuite.GetBaseURL()+"/api/subscribe",
		bytes.NewBuffer(jsonBody))
	assert.NoError(s.T(), err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
}

func (s *SimpleSubscriptionEndpointsTestSuite) TestConfirmUnsubscribeFlow() {
	client := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

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

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		s.testSuite.GetBaseURL()+"/api/confirm/test-confirmation-token", nil)
	assert.NoError(s.T(), err)

	resp, err := client.Do(req)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusFound, resp.StatusCode)
	location := resp.Header.Get("Location")
	assert.Contains(s.T(), location, "message_type=success")

	var updatedSubscription models.Subscription
	err = s.testSuite.DB.Where("id = ?", subscription.ID).First(&updatedSubscription).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "active", updatedSubscription.Status)

	req, err = http.NewRequestWithContext(ctx, http.MethodGet,
		s.testSuite.GetBaseURL()+"/api/unsubscribe/test-unsubscribe-token", nil)
	assert.NoError(s.T(), err)

	resp, err = client.Do(req)
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	assert.Equal(s.T(), http.StatusFound, resp.StatusCode)
	location = resp.Header.Get("Location")
	assert.Contains(s.T(), location, "message_type=success")

	err = s.testSuite.DB.Where("id = ?", subscription.ID).First(&updatedSubscription).Error
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "cancelled", updatedSubscription.Status)
}

func TestSimpleSubscriptionEndpointsTestSuite(t *testing.T) {
	suite.Run(t, new(SimpleSubscriptionEndpointsTestSuite))
}
