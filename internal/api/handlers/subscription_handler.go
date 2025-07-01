package handlers

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ihorlenko/weather_notifier/internal/interfaces"
)

type SubscriptionHandler struct {
	orchestrator        interfaces.SubscriptionOrchestrator
	subscriptionService interfaces.SubscriptionService
}

type SubscribeRequest struct {
	Email     string `json:"email" binding:"required,email"`
	City      string `json:"city" binding:"required"`
	Frequency string `json:"frequency" binding:"required,oneof=hourly daily"`
}

func NewSubscriptionHandler(
	orchestrator interfaces.SubscriptionOrchestrator,
	subscriptionService interfaces.SubscriptionService,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		orchestrator:        orchestrator,
		subscriptionService: subscriptionService,
	}
}

// Subscribe godoc
// @Summary      Subscribe for weather updates
// @Description  Subscribes a given email for weather updates
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request body SubscribeRequest true "Subscription data"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscribe [post]
func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := h.orchestrator.ProcessSubscription(ctx, req.Email, req.City, req.Frequency)
	if err != nil {
		h.handleSubscriptionError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Please check your email to confirm the subscription",
	})
}

// Confirm godoc
// @Summary      Confirm subscription
// @Description  Confirm subscription via token from letter
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        token path string true "Confirmation token"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /confirm/{token} [get]
func (h *SubscriptionHandler) Confirm(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		redirectURL := "/?message_type=error&message=" + url.QueryEscape("Token is required")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	err := h.subscriptionService.ConfirmSubscription(token)
	if err != nil {
		redirectURL := "/?message_type=error&message=" + url.QueryEscape(err.Error())
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	message := "Your subscription has been successfully confirmed!"
	redirectURL := "/?message_type=success&message=" + url.QueryEscape(message)
	c.Redirect(http.StatusFound, redirectURL)
}

// Unsubscribe godoc
// @Summary      Unsubscribe from updates
// @Description  Unsubscribes user from weather updates
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        token path string true "Unsubscribe token"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /unsubscribe/{token} [get]
func (h *SubscriptionHandler) Unsubscribe(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		redirectURL := "/?message_type=error&message=" + url.QueryEscape("Token is required")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	err := h.subscriptionService.Unsubscribe(token)
	if err != nil {
		redirectURL := "/?message_type=error&message=" + url.QueryEscape(err.Error())
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	message := "You have successfully unsubscribed from weather updates"
	redirectURL := "/?message_type=success&message=" + url.QueryEscape(message)
	c.Redirect(http.StatusFound, redirectURL)
}

func (h *SubscriptionHandler) handleSubscriptionError(c *gin.Context, err error) {
	errorMsg := err.Error()

	switch {
	case contains(errorMsg, "city validation failed"):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city or weather service unavailable"})
	case contains(errorMsg, "subscription creation failed"):
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create subscription"})
	case contains(errorMsg, "confirmation email failed"):
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Subscription created but email failed to send"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
	}
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr
}
